package hash

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// ChunkedHasher is a hash.Hash wrapper that will hash data in chunks, similar to how a torrent client hashes data.
//
// Call [Flush] to ensure that any misaligned data at the tail is hashed.
// Call [Close] to release resources.
type ChunkedHasher struct {
	ht              Type
	h               hash.Hash
	inBuf           []byte
	outBuf          []byte
	outReadIdx      int
	outWriteIdx     int
	chunkSize       int
	hashedChunkSize int
	flushPending    bool
	mu              sync.Mutex
}

// NewChunkedHasher will return a new ChunkedHasher that will hash data in chunks of the given size.
//
// If chunkSize is -1, the hasher will use the block size of the hash function;
// Use [IngressChunkSize] to get the resulting chunk size.
func NewChunkedHasher(ht Type, chunkSize int) *ChunkedHasher {
	h := getHasher(ht)
	buf := bufPool.Get().([]byte)
	buf = buf[:0]
	if chunkSize == -1 {
		chunkSize = h.BlockSize()
	}
	if cap(buf) < chunkSize {
		buf = append(buf, make([]byte, chunkSize-cap(buf))...)
	}

	ch := &ChunkedHasher{
		ht:         ht,
		h:          h,
		inBuf:      buf,
		chunkSize:  chunkSize,
		outBuf:     bufPool.Get().([]byte),
		outReadIdx: 0,
	}

	dummyBuf := bufPool.Get().([]byte)
	if cap(dummyBuf) < chunkSize {
		dummyBuf = append(dummyBuf, make([]byte, chunkSize-cap(dummyBuf))...)
	}
	if len(dummyBuf) < chunkSize {
		dummyBuf = dummyBuf[:chunkSize]
	}
	for i := 0; i < chunkSize; i++ {
		dummyBuf[i] = 5
	}
	ch.h.Write(dummyBuf)
	ch.hashedChunkSize = ch.h.Size()
	ch.h.Reset()
	bufPool.Put(dummyBuf)
	return ch
}

func (c *ChunkedHasher) Close() {
	c.mu.Lock()
	putHasher(c.ht, c.h)
	bufPool.Put(c.inBuf)
	bufPool.Put(c.outBuf)
	c.inBuf = nil
	c.outBuf = nil
	c.h = nil
	c.mu.Unlock()
}

var ErrHasherClosed = errors.New("hasher is closed")

// spin manages outWriteIdx and outReadIdx as ring buffer style indices.
// caller MUST hold the lock.
func (c *ChunkedHasher) spin(read, written int, flushing bool) {
	c.outReadIdx += read
	c.outWriteIdx += written
	if !flushing && (c.outReadIdx != 0 && c.outReadIdx%c.hashedChunkSize != 0) {
		spew.Dump(c)
		panic("someone forgot to manage the indexes")
	}
	if flushing {
		c.outWriteIdx = 0
	}
	if c.outReadIdx >= len(c.outBuf) {
		c.outReadIdx = 0
	}
	if flushing {
		return
	}
	if c.outWriteIdx >= len(c.outBuf) && c.outReadIdx == 0 {
		c.outWriteIdx = 0
	}
	if c.outWriteIdx >= len(c.outBuf) {
		c.outBuf = append(c.outBuf, make([]byte, c.hashedChunkSize)...)
	}
}

// sum will hash the current chunk of data, caller MUST hold the lock.
func (c *ChunkedHasher) sum(flushing bool) error {
	if c.inBuf == nil || c.h == nil || c.outBuf == nil {
		return ErrHasherClosed
	}
	if len(c.inBuf) == 0 {
		return nil
	}

	inSize := len(c.inBuf)

	switch {
	case inSize < c.chunkSize && !flushing:
		return nil
	case inSize < c.chunkSize && flushing:
		c.h.Write(c.inBuf[:inSize])
		c.spin(0, c.h.Size(), flushing)
		c.outBuf = c.h.Sum(c.outBuf)
		c.h.Reset()
		c.inBuf = c.inBuf[:0]
		return nil
	case inSize > c.chunkSize:
		for len(c.inBuf) >= c.chunkSize {
			c.h.Write(c.inBuf[:c.chunkSize])
			c.spin(0, c.h.Size(), flushing)
			c.outBuf = c.h.Sum(c.outBuf)
			c.h.Reset()
			c.inBuf = c.inBuf[c.chunkSize:]
		}
		if len(c.inBuf) > 0 && (flushing || len(c.inBuf) == c.chunkSize) {
			c.h.Write(c.inBuf)
			c.spin(0, c.h.Size(), flushing)
			c.outBuf = c.h.Sum(c.outBuf)
			c.h.Reset()
			c.inBuf = c.inBuf[:0]
		}
		return nil
	case inSize == c.chunkSize:
		c.h.Write(c.inBuf)
		c.spin(0, c.h.Size(), flushing)
		c.outBuf = c.h.Sum(c.outBuf)
		c.h.Reset()
		c.inBuf = c.inBuf[:0]
		return nil
	default:
		panic("unreachable, or so we hope")
	}
}

// Flush will ensure that any misaligned data at the tail is hashed.
// Must be called before calling [Next] or [Write] again, but after all data has been written.
//
// After calling Flush, the index pointer will be reset to the beginning of the output buffer.
// This means that if you don't align this call with the end of your dataset, you will lose data.
func (c *ChunkedHasher) Flush() error {
	c.mu.Lock()
	if err := c.sum(true); err != nil {
		c.mu.Unlock()
		return err
	}
	c.mu.Unlock()
	return nil
}

func (c *ChunkedHasher) IngressChunkSize() int {
	return c.chunkSize
}

func (c *ChunkedHasher) EgressChunkSize() int {
	return c.hashedChunkSize
}

func (c *ChunkedHasher) Write(p []byte) (n int, err error) {
	n = len(p)
	if len(p) == 0 {
		return
	}
	c.mu.Lock()
	if c.inBuf == nil || c.h == nil || c.outBuf == nil {
		c.mu.Unlock()
		return 0, ErrHasherClosed
	}
	c.inBuf = append(c.inBuf, p...)
	if len(c.inBuf) >= c.chunkSize {
		if err = c.sum(false); err != nil {
			c.mu.Unlock()
			return 0, err
		}
	}
	c.mu.Unlock()
	return
}

// Next writes the hash of the next chunk of data to dst.
//
// Returns the resulting slice, the number of bytes written, and any error that occurred.
// Note that the number of bytes will always be the chunk size, unless an error occurred.
//
//   - If dst is too small, it will be resized to fit the chunk size.
//   - If dst is too large, it will be truncated to the chunk size.
//   - Any existing data in dst will be overwritten.
//   - If there is no data to hash, it will return an error.
func (c *ChunkedHasher) Next(dst []byte) ([]byte, int, error) {
	c.mu.Lock()
	if c.inBuf == nil || c.h == nil || c.outBuf == nil {
		c.mu.Unlock()
		return dst, 0, ErrHasherClosed
	}
	if err := c.sum(false); err != nil {
		c.mu.Unlock()
		return dst, 0, err
	}
	if len(c.outBuf) == 0 {
		c.mu.Unlock()
		return dst, 0, fmt.Errorf("%w: no data to hash", io.EOF)
	}
	if c.outReadIdx >= len(c.outBuf) {
		panic("someone forgot to move the index pointer")
	}
	if cap(dst) < len(dst)+c.hashedChunkSize {
		dst = append(dst, make([]byte, c.hashedChunkSize)...)
	}
	dst = dst[:c.hashedChunkSize]
	n := copy(dst, c.outBuf[c.outReadIdx:c.outReadIdx+c.hashedChunkSize])
	if n != c.hashedChunkSize {
		c.mu.Unlock()
		spew.Dump(c)
		spew.Dump(dst)
		// panic justification: handling the index pointer after a failure like this is complex. FUBAR.
		// caller just needs to get it together and not fuck this up; else take a panic to the face.
		panic(fmt.Errorf("%w: expected to copy %d bytes, but only copied %d", io.ErrShortWrite, c.hashedChunkSize, n))
	}
	c.spin(n, 0, false)
	c.mu.Unlock()
	return dst, n, nil
}
