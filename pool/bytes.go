package pool

import (
	"bytes"
	"io"
	"sync"
)

type BufferFactory struct {
	pool *sync.Pool
}

func NewBufferFactory() BufferFactory {
	return BufferFactory{
		pool: &sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		},
	}
}

func (cf BufferFactory) Put(buf *Buffer) error {
	var err = ErrBufferReturned
	buf.o.Do(func() {
		_ = buf.Reset()
		cf.pool.Put(buf.Buffer)
		buf.Buffer = nil
		err = nil
	})
	return err
}

func (cf BufferFactory) MustPut(buf *Buffer) {
	if err := cf.Put(buf); err != nil {
		panic(err)
	}
}

func (cf BufferFactory) Get() *Buffer {
	return &Buffer{
		cf.pool.Get().(*bytes.Buffer),
		&sync.Once{},
	}
}

type Buffer struct {
	*bytes.Buffer
	o *sync.Once
}

func (c Buffer) Bytes() []byte {
	if c.Buffer == nil {
		return nil
	}
	return c.Buffer.Bytes()
}

func (c Buffer) MustBytes() []byte {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.Bytes()
}

func (c Buffer) String() string {
	if c.Buffer == nil {
		return ""
	}
	return c.Buffer.String()
}

func (c Buffer) MustString() string {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.String()
}

func (c Buffer) Reset() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Reset()
	return nil
}

func (c Buffer) MustReset() {
	if err := c.Reset(); err != nil {
		panic(err)
	}
	c.Buffer.Reset()
}

func (c Buffer) Len() int {
	if c.Buffer == nil {
		return 0
	}
	return c.Buffer.Len()
}

func (c Buffer) MustLen() int {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.Len()
}

func (c Buffer) Write(p []byte) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.Write(p)
}

func (c Buffer) MustWrite(p []byte) {
	if _, err := c.Write(p); err != nil {
		panic(err)
	}
}

func (c Buffer) WriteRune(r rune) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteRune(r)
}

func (c Buffer) MustWriteRune(r rune) {
	if _, err := c.WriteRune(r); err != nil {
		panic(err)
	}
}

func (c Buffer) WriteByte(cyte byte) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.WriteByte(cyte)
}

func (c Buffer) MustWriteByte(cyte byte) {
	if err := c.WriteByte(cyte); err != nil {
		panic(err)
	}
}

func (c Buffer) WriteString(str string) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteString(str)
}

func (c Buffer) Grow(n int) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Grow(n)
	return nil
}

func (c Buffer) Cap() int {
	if c.Buffer == nil {
		return 0
	}
	return c.Buffer.Cap()
}

func (c Buffer) Truncate(n int) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Truncate(n)
	return nil
}

func (c Buffer) MustTruncate(n int) {
	if err := c.Truncate(n); err != nil {
		panic(err)
	}
}

func (c Buffer) ReadFrom(r io.Reader) (int64, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.ReadFrom(r)
}

func (c Buffer) MustReadFrom(r io.Reader) {
	if _, err := c.ReadFrom(r); err != nil {
		panic(err)
	}
}

func (c Buffer) WriteTo(w io.Writer) (int64, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteTo(w)
}

func (c Buffer) MustWriteTo(w io.Writer) {
	if _, err := c.WriteTo(w); err != nil {
		panic(err)
	}
}

func (c Buffer) Read(p []byte) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.Read(p)
}

func (c Buffer) ReadByte() (byte, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.ReadByte()
}

func (c Buffer) ReadRune() (rune, int, error) {
	if c.Buffer == nil {
		return 0, 0, ErrBufferReturned
	}
	return c.Buffer.ReadRune()
}

func (c Buffer) UnreadByte() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.UnreadByte()
}

func (c Buffer) UnreadRune() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.UnreadRune()
}

func (c Buffer) ReadBytes(delim byte) ([]byte, error) {
	if c.Buffer == nil {
		return nil, ErrBufferReturned
	}
	return c.Buffer.ReadBytes(delim)
}

func (c Buffer) Next(n int) []byte {
	if c.Buffer == nil {
		return nil
	}
	return c.Buffer.Next(n)
}
