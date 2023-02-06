package pool

import (
	"bytes"
	"io"
	"sync"
)

// BufferFactory is a factory for creating and reusing bytes.Buffers.
// BufferFactory tries to be safer than using a sync.Pool directly by ensuring that the buffer is not returned twice.
type BufferFactory struct {
	pool *sync.Pool
}

// NewBufferFactory creates a new BufferFactory that creates new buffers on demand.
func NewBufferFactory() BufferFactory {
	return BufferFactory{
		pool: &sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		},
	}
}

// NewSizedBufferFactory creates a new BufferFactory that creates new buffers of the given size on demand.
func NewSizedBufferFactory(size int) BufferFactory {
	return BufferFactory{
		pool: &sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, size)) },
		},
	}
}

// Put returns the buffer to the pool. It returns an error if the buffer has already been returned to the pool.
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

// MustPut is the same as Put but panics if the buffer has already been returned to the pool.
func (cf BufferFactory) MustPut(buf *Buffer) {
	if err := cf.Put(buf); err != nil {
		panic(err)
	}
}

// Get returns a buffer from the pool.
func (cf BufferFactory) Get() *Buffer {
	return &Buffer{
		cf.pool.Get().(*bytes.Buffer),
		&sync.Once{},
	}
}

// Buffer is a wrapper around bytes.Buffer that can only be returned to a pool once.
type Buffer struct {
	*bytes.Buffer
	o *sync.Once
}

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
//
// *This is from the bytes.Buffer docs.*
func (c Buffer) Bytes() []byte {
	if c.Buffer == nil {
		return nil
	}
	return c.Buffer.Bytes()
}

// MustBytes is the same as Bytes but panics if the buffer has already been returned to the pool.
func (c Buffer) MustBytes() []byte {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.Bytes()
}

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
//
// *This is from the bytes.Buffer docs.*
func (c Buffer) String() string {
	if c.Buffer == nil {
		return ""
	}
	return c.Buffer.String()
}

// MustString is the same as String but panics if the buffer has already been returned to the pool.
func (c Buffer) MustString() string {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.String()
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) Reset() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Reset()
	return nil
}

// MustReset is the same as Reset but panics if the buffer has already been returned to the pool.
func (c Buffer) MustReset() {
	if err := c.Reset(); err != nil {
		panic(err)
	}
	c.Buffer.Reset()
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns 0 if the buffer has already been returned to the pool.
func (c Buffer) Len() int {
	if c.Buffer == nil {
		return 0
	}
	return c.Buffer.Len()
}

// MustLen is the same as Len but panics if the buffer has already been returned to the pool.
func (c Buffer) MustLen() int {
	if c.Buffer == nil {
		panic(ErrBufferReturned)
	}
	return c.Buffer.Len()
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) Write(p []byte) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.Write(p)
}

// MustWrite is the same as Write but panics if the buffer has already been returned to the pool.
func (c Buffer) MustWrite(p []byte) {
	if _, err := c.Write(p); err != nil {
		panic(err)
	}
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) WriteRune(r rune) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteRune(r)
}

// MustWriteRune is the same as WriteRune but panics if the buffer has already been returned to the pool.
func (c Buffer) MustWriteRune(r rune) {
	if _, err := c.WriteRune(r); err != nil {
		panic(err)
	}
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) WriteByte(cyte byte) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.WriteByte(cyte)
}

// MustWriteByte is the same as WriteByte but panics if the buffer has already been returned to the pool.
func (c Buffer) MustWriteByte(cyte byte) {
	if err := c.WriteByte(cyte); err != nil {
		panic(err)
	}
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) WriteString(str string) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteString(str)
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for another n bytes. After Grow(n), at least n bytes can be written to the buffer without another allocation. If n is negative, Grow will panic. If the buffer can't grow it will panic with ErrTooLarge.
// If the buffer has already been returned to the pool, Grow will return ErrBufferReturned.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) Grow(n int) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Grow(n)
	return nil
}

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
//
// *This is from the bytes.Buffer docs.*
// If the buffer has already been returned to the pool, Cap will return 0.
func (c Buffer) Cap() int {
	if c.Buffer == nil {
		return 0
	}
	return c.Buffer.Cap()
}

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) Truncate(n int) error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	c.Buffer.Truncate(n)
	return nil
}

// MustTruncate is the same as Truncate but panics if the buffer has already been returned to the pool.
func (c Buffer) MustTruncate(n int) {
	if err := c.Truncate(n); err != nil {
		panic(err)
	}
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) ReadFrom(r io.Reader) (int64, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.ReadFrom(r)
}

// MustReadFrom is the same as ReadFrom but panics if the buffer has already been returned to the pool.
func (c Buffer) MustReadFrom(r io.Reader) {
	if _, err := c.ReadFrom(r); err != nil {
		panic(err)
	}
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) WriteTo(w io.Writer) (int64, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.WriteTo(w)
}

// MustWriteTo is the same as WriteTo but panics if the buffer has already been returned to the pool.
func (c Buffer) MustWriteTo(w io.Writer) {
	if _, err := c.WriteTo(w); err != nil {
		panic(err)
	}
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) Read(p []byte) (int, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.Read(p)
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) ReadByte() (byte, error) {
	if c.Buffer == nil {
		return 0, ErrBufferReturned
	}
	return c.Buffer.ReadByte()
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) ReadRune() (rune, int, error) {
	if c.Buffer == nil {
		return 0, 0, ErrBufferReturned
	}
	return c.Buffer.ReadRune()
}

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) UnreadByte() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.UnreadByte()
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a successful ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) UnreadRune() error {
	if c.Buffer == nil {
		return ErrBufferReturned
	}
	return c.Buffer.UnreadRune()
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns an error if the buffer has already been returned to the pool.
func (c Buffer) ReadBytes(delim byte) ([]byte, error) {
	if c.Buffer == nil {
		return nil, ErrBufferReturned
	}
	return c.Buffer.ReadBytes(delim)
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
//
// *This is from the bytes.Buffer docs.*
// This wrapper returns nil if the buffer has already been returned to the pool.
func (c Buffer) Next(n int) []byte {
	if c.Buffer == nil {
		return nil
	}
	return c.Buffer.Next(n)
}
