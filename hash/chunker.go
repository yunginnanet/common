package hash

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"slices"
	"sync"
)

// ChunkedHasher is a hash.Hash wrapper that will hash data in chunks, similar to how a torrent client hashes data.
//
// Call [Flush] to ensure that any misaligned data at the tail is hashed.
// Call [Close] to release resources.
type ChunkedHasher struct {
	ht              Type
	h               hash.Hash
	inBuf           []byte
	outBufs         [][]byte
	chunkSize       int
	hashedChunkSize int
	mu              sync.Mutex
}

// NewChunkedHasher will return a new ChunkedHasher that will hash data in chunks of the given size.
//
// If chunkSize is -1, the hasher will use the block size of the hash function;
// Use [IngressChunkSize] to get the resulting chunk size.
func NewChunkedHasher(ht Type, chunkSize int) *ChunkedHasher {
	h := getHasher(ht)
	if chunkSize == -1 {
		chunkSize = h.BlockSize()
	}

	newInBuf := bufPool.Get().([]byte)
	clear(newInBuf)
	if cap(newInBuf) < chunkSize {
		newInBuf = slices.Grow(newInBuf, chunkSize)
	}
	newInBuf = newInBuf[:chunkSize]

	ch := &ChunkedHasher{
		ht:        ht,
		h:         h,
		inBuf:     newInBuf,
		chunkSize: chunkSize,
		outBufs:   make([][]byte, 0, 1),
	}

	newOutBuf := bufPool.Get().([]byte)
	clear(newOutBuf)
	if cap(newOutBuf) < chunkSize {
		newOutBuf = slices.Grow(newOutBuf, chunkSize)
	}
	newOutBuf = newOutBuf[:chunkSize]

	for i := 0; i < chunkSize; i++ {
		newOutBuf[i] = 5
	}

	ch.h.Write(newOutBuf[:chunkSize])
	_ = ch.h.Sum(nil)
	ch.hashedChunkSize = ch.h.Size()
	ch.h.Reset()
	clear(newOutBuf)

	ch.outBufs = append(ch.outBufs, newInBuf[:chunkSize])

	return ch
}

func (c *ChunkedHasher) Close() {
	c.mu.Lock()
	putHasher(c.ht, c.h)
	bufPool.Put(c.inBuf)
	for _, b := range c.outBufs {
		clear(b)
		bufPool.Put(b)
		b = nil
	}
	c.inBuf = nil
	c.h = nil
	c.mu.Unlock()
}

var (
	ErrHasherClosed = errors.New("hasher is closed")
)

var ErrMrHopefulUnreachable = errors.New("unreachable, or so we hope")

// sum will hash the current chunk of data, caller MUST hold the lock.
func (c *ChunkedHasher) sum(flushing bool) error {
	if c.inBuf == nil || c.h == nil || c.outBufs == nil {
		return ErrHasherClosed
	}

	if len(c.inBuf) == 0 {
		return nil
	}

	inSize := len(c.inBuf)
	if inSize == 0 && !flushing {
		return nil
	}

	latestBuf := len(c.outBufs) - 1

	if flushing && len(c.inBuf) > 0 {
		switch {
		case len(c.inBuf) > c.chunkSize:
			panic(ErrMrHopefulUnreachable)
		case len(c.inBuf) < c.chunkSize:
			c.h.Reset()
			flushed := 0
			for {
				n, e := c.h.Write(c.inBuf[:c.chunkSize-len(c.inBuf)])
				if e != nil {
					return e
				}
				if n == 0 {
					break
				}
				c.inBuf = c.inBuf[n:]
				flushed += n
				if flushed == c.chunkSize || len(c.inBuf) == 0 || c.h.Size() == c.hashedChunkSize {
					break
				}
			}
		}
		switch {
		case len(c.outBufs[latestBuf]) == c.hashedChunkSize:
			newBuf := bufPool.Get().([]byte)
			clear(newBuf)
			if cap(newBuf) < c.chunkSize {
				newBuf = slices.Grow(newBuf, c.chunkSize-cap(newBuf))
			}
			newBuf = newBuf[:c.h.Size()]
			newBuf = c.h.Sum(newBuf)
			c.outBufs = append(c.outBufs, newBuf[:c.h.Size()])
			latestBuf++
		case len(c.outBufs[latestBuf]) == 0:
			c.outBufs[latestBuf] = c.h.Sum(c.outBufs[latestBuf])
		case len(c.outBufs[latestBuf]) < c.hashedChunkSize && len(c.outBufs[latestBuf]) > 0:

			panic(ErrMrHopefulUnreachable)
		}
	}

	if !flushing && len(c.inBuf) == c.chunkSize {
		c.h.Reset()
		n, e := c.h.Write(c.inBuf)
		if e != nil {
			return e
		}
		if n != c.chunkSize {
			return fmt.Errorf("%w: expected to write %d bytes, but only wrote %d", io.ErrShortWrite, c.chunkSize, n)
		}
		newBuf := bufPool.Get().([]byte)
		clear(newBuf)
		if cap(newBuf) < c.chunkSize {
			newBuf = slices.Grow(newBuf, c.chunkSize-cap(newBuf))
		}
		newBuf = newBuf[:c.h.Size()]
		newBuf = c.h.Sum(newBuf)
		c.outBufs = append(c.outBufs, newBuf[:c.h.Size()])
		latestBuf++
		c.h.Reset()
		c.inBuf = c.inBuf[:n]
	}

	return nil
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
	if n = len(p); n == 0 {
		return
	}
	c.mu.Lock()
	if c.inBuf == nil || c.h == nil || c.outBufs == nil {
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
	if c.inBuf == nil || c.h == nil || c.outBufs == nil {
		c.mu.Unlock()
		return dst, 0, ErrHasherClosed
	}

	c.outBufs = slices.Clip(c.outBufs)

	if len(c.outBufs) == 0 || len(c.outBufs) == 1 && len(c.outBufs[0]) == 0 {
		c.mu.Unlock()
		println("EOF")
		return dst, 0, io.EOF
	}

	if dst == nil {
		dst = bufPool.Get().([]byte)
		clear(dst)
		dst = dst[:c.hashedChunkSize]
	}

	if cap(dst) < c.hashedChunkSize {
		dst = slices.Grow(dst, c.hashedChunkSize-cap(dst))
	}
	dst = dst[:c.hashedChunkSize]
	n := copy(dst, c.outBufs[len(c.outBufs)-1])
	if n == len(c.outBufs[len(c.outBufs)-1]) {
		c.outBufs = slices.Delete(c.outBufs, len(c.outBufs)-1, len(c.outBufs))
		c.outBufs = slices.Clip(c.outBufs)
	} else {
		c.outBufs[len(c.outBufs)-1] = c.outBufs[len(c.outBufs)-1][n:]
		c.mu.Unlock()
		return dst, n, io.ErrShortBuffer
	}

	c.mu.Unlock()
	return dst, n, nil
}
