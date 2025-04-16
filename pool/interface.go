package pool

import (
	"bytes"
	"io"
	"sync"
)

type Pool[T any] interface {
	Get() T
	Put(T)
}

// ByteBuffer is satisfied by [*pool.Buffer] and [*bytes.Buffer].
//
// Note, we can't include Reset and Grow as our implementations return an error while a [*bytes.Buffer] does not.
type ByteBuffer interface {
	Len() int
	Cap() int
	Bytes() []byte
	WriteRune(rune) (int, error)
	WriteString(string) (int, error)

	io.WriterTo
	io.ByteWriter
	io.ReadWriter

	io.ReaderFrom
	io.ByteReader
	io.RuneReader
}

type WithPutError[T any] interface {
	Get() T
	Put(T) error
}

func (b BufferFactoryInterfaceCompat) Put(buf *Buffer) {
	_ = b.BufferFactory.Put(buf)
}

type BufferFactoryInterfaceCompat struct {
	BufferFactory
}

type BufferFactoryByteBufferCompat struct {
	BufferFactory
}

func (bf BufferFactoryByteBufferCompat) Put(buf ByteBuffer) {
	if b, ok := buf.(*Buffer); ok {
		err := bf.BufferFactory.Put(b)
		if err != nil {
			panic(err)
		}
		return
	}
	if b, ok := buf.(*bytes.Buffer); ok {
		newB := &Buffer{
			o:      &sync.Once{},
			Buffer: b,
		}
		_ = bf.BufferFactory.Put(newB)
		return
	}
	// unfortunately this compatibility shim cannot be used with any other types implementing ByteBuffer
	// this is because we can't wrap them in a *Buffer
	panic("invalid type, need *pool.Buffer or *bytes.Buffer")
}

func (bf BufferFactoryByteBufferCompat) Get() ByteBuffer {
	b := bf.BufferFactory.Get()
	return ByteBuffer(b)
}

var (
	_ ByteBuffer       = (*Buffer)(nil)
	_ ByteBuffer       = (*bytes.Buffer)(nil)
	_ Pool[ByteBuffer] = (*BufferFactoryByteBufferCompat)(nil)
)
