# pool



```go
var ErrBufferReturned = errors.New("buffer already returned")
```

#### type Buffer

```go
type Buffer struct {
	*bytes.Buffer
}
```

Buffer is a wrapper around bytes.Buffer that can only be returned to a pool
once.

#### func (Buffer) Bytes

```go
func (c Buffer) Bytes() []byte
```
Bytes returns a slice of length b.Len() holding the unread portion of the
buffer. The slice is valid for use only until the next buffer modification (that
is, only until the next call to a method like Read, Write, Reset, or Truncate).
The slice aliases the buffer content at least until the next buffer
modification, so immediate changes to the slice will affect the result of future
reads.

*This is from the bytes.Buffer docs.*

#### func (Buffer) Cap

```go
func (c Buffer) Cap() int
```
Cap returns the capacity of the buffer's underlying byte slice, that is, the
total space allocated for the buffer's data.

*This is from the bytes.Buffer docs.* If the buffer has already been returned to
the pool, Cap will return 0.

#### func (Buffer) Close

```go
func (c Buffer) Close() error
```
Close implements io.Closer. It returns the buffer to the pool. This

#### func (Buffer) Grow

```go
func (c Buffer) Grow(n int) error
```
Grow grows the buffer's capacity, if necessary, to guarantee space for another n
bytes. After Grow(n), at least n bytes can be written to the buffer without
another allocation. If n is negative, Grow will panic. If the buffer can't grow
it will panic with ErrTooLarge.

If the buffer has already been returned to the pool, Grow will return
ErrBufferReturned.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) IsClosed

```go
func (c Buffer) IsClosed() bool
```
IsClosed returns true if the buffer has been returned to the pool.

#### func (Buffer) Len

```go
func (c Buffer) Len() int
```
Len returns the number of bytes of the unread portion of the buffer; b.Len() ==
len(b.Bytes()).

*This is from the bytes.Buffer docs.* This wrapper returns 0 if the buffer has
already been returned to the pool.

#### func (Buffer) MustBytes

```go
func (c Buffer) MustBytes() []byte
```
MustBytes is the same as Bytes but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustLen

```go
func (c Buffer) MustLen() int
```
MustLen is the same as Len but panics if the buffer has already been returned to
the pool.

#### func (Buffer) MustReadFrom

```go
func (c Buffer) MustReadFrom(r io.Reader)
```
MustReadFrom is the same as ReadFrom but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustReset

```go
func (c Buffer) MustReset()
```
MustReset is the same as Reset but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustString

```go
func (c Buffer) MustString() string
```
MustString is the same as String but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustTruncate

```go
func (c Buffer) MustTruncate(n int)
```
MustTruncate is the same as Truncate but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustWrite

```go
func (c Buffer) MustWrite(p []byte)
```
MustWrite is the same as Write but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustWriteByte

```go
func (c Buffer) MustWriteByte(cyte byte)
```
MustWriteByte is the same as WriteByte but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustWriteRune

```go
func (c Buffer) MustWriteRune(r rune)
```
MustWriteRune is the same as WriteRune but panics if the buffer has already been
returned to the pool.

#### func (Buffer) MustWriteTo

```go
func (c Buffer) MustWriteTo(w io.Writer)
```
MustWriteTo is the same as WriteTo but panics if the buffer has already been
returned to the pool.

#### func (Buffer) Next

```go
func (c Buffer) Next(n int) []byte
```
Next returns a slice containing the next n bytes from the buffer, advancing the
buffer as if the bytes had been returned by Read. If there are fewer than n
bytes in the buffer, Next returns the entire buffer. The slice is only valid
until the next call to a read or write method.

*This is from the bytes.Buffer docs.* This wrapper returns nil if the buffer has
already been returned to the pool.

#### func (Buffer) Read

```go
func (c Buffer) Read(p []byte) (int, error)
```
Read reads the next len(p) bytes from the buffer or until the buffer is drained.
The return value n is the number of bytes read. If the buffer has no data to
return, err is io.EOF (unless len(p) is zero); otherwise it is nil.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) ReadByte

```go
func (c Buffer) ReadByte() (byte, error)
```
ReadByte reads and returns the next byte from the buffer. If no byte is
available, it returns error io.EOF.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) ReadBytes

```go
func (c Buffer) ReadBytes(delim byte) ([]byte, error)
```
ReadBytes reads until the first occurrence of delim in the input, returning a
slice containing the data up to and including the delimiter. If ReadBytes
encounters an error before finding a delimiter, it returns the data read before
the error and the error itself (often io.EOF). ReadBytes returns err != nil if
and only if the returned data does not end in delim.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) ReadFrom

```go
func (c Buffer) ReadFrom(r io.Reader) (int64, error)
```
ReadFrom reads data from r until EOF and appends it to the buffer, growing the
buffer as needed. The return value n is the number of bytes read. Any error
except io.EOF encountered during the read is also returned. If the buffer
becomes too large, ReadFrom will panic with ErrTooLarge.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) ReadRune

```go
func (c Buffer) ReadRune() (rune, int, error)
```
ReadRune reads and returns the next UTF-8-encoded Unicode code point from the
buffer. If no bytes are available, the error returned is io.EOF. If the bytes
are an erroneous UTF-8 encoding, it consumes one byte and returns U+FFFD, 1.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) Reset

```go
func (c Buffer) Reset() error
```
Reset resets the buffer to be empty, but it retains the underlying storage for
use by future writes. Reset is the same as Truncate(0).

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) String

```go
func (c Buffer) String() string
```
String returns the contents of the unread portion of the buffer as a string. If
the Buffer is a nil pointer, it returns "<nil>".

To build strings more efficiently, see the strings.Builder type.

*This is from the bytes.Buffer docs.*

#### func (Buffer) Truncate

```go
func (c Buffer) Truncate(n int) error
```
Truncate discards all but the first n unread bytes from the buffer but continues
to use the same allocated storage. It panics if n is negative or greater than
the length of the buffer.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) UnreadByte

```go
func (c Buffer) UnreadByte() error
```
UnreadByte unreads the last byte returned by the most recent successful read
operation that read at least one byte. If a write has happened since the last
read, if the last read returned an error, or if the read read zero bytes,
UnreadByte returns an error.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) UnreadRune

```go
func (c Buffer) UnreadRune() error
```
UnreadRune unreads the last rune returned by ReadRune. If the most recent read
or write operation on the buffer was not a successful ReadRune, UnreadRune
returns an error. (In this regard it is stricter than UnreadByte, which will
unread the last byte from any read operation.)

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) WithParent

```go
func (c Buffer) WithParent(p *BufferFactory) *Buffer
```
WithParent sets the parent of the buffer. This is useful for chaining factories,
and for facilitating in-line buffer return with functions like Buffer.Close().
Be mindful, however, that this adds a bit of overhead.

#### func (Buffer) Write

```go
func (c Buffer) Write(p []byte) (int, error)
```
Write appends the contents of p to the buffer, growing the buffer as needed. The
return value n is the length of p; err is always nil. If the buffer becomes too
large, Write will panic with ErrTooLarge.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) WriteByte

```go
func (c Buffer) WriteByte(cyte byte) error
```
WriteByte appends the byte c to the buffer, growing the buffer as needed. The
returned error is always nil, but is included to match bufio.Writer's WriteByte.
If the buffer becomes too large, WriteByte will panic with ErrTooLarge.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) WriteRune

```go
func (c Buffer) WriteRune(r rune) (int, error)
```
WriteRune appends the UTF-8 encoding of Unicode code point r to the buffer,
returning its length and an error, which is always nil but is included to match
bufio.Writer's WriteRune. The buffer is grown as needed; if it becomes too
large, WriteRune will panic with ErrTooLarge.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) WriteString

```go
func (c Buffer) WriteString(str string) (int, error)
```
WriteString appends the contents of s to the buffer, growing the buffer as
needed. The return value n is the length of s; err is always nil. If the buffer
becomes too large, WriteString will panic with ErrTooLarge.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### func (Buffer) WriteTo

```go
func (c Buffer) WriteTo(w io.Writer) (int64, error)
```
WriteTo writes data to w until the buffer is drained or an error occurs. The
return value n is the number of bytes written; it always fits into an int, but
it is int64 to match the io.WriterTo interface. Any error encountered during the
write is also returned.

*This is from the bytes.Buffer docs.* This wrapper returns an error if the
buffer has already been returned to the pool.

#### type BufferFactory

```go
type BufferFactory struct {
}
```

BufferFactory is a factory for creating and reusing bytes.Buffers. BufferFactory
tries to be safer than using a sync.Pool directly by ensuring that the buffer is
not returned twice.

#### func  NewBufferFactory

```go
func NewBufferFactory() BufferFactory
```
NewBufferFactory creates a new BufferFactory that creates new buffers on demand.

#### func  NewSizedBufferFactory

```go
func NewSizedBufferFactory(size int) BufferFactory
```
NewSizedBufferFactory creates a new BufferFactory that creates new buffers of
the given size on demand.

#### func (BufferFactory) Get

```go
func (cf BufferFactory) Get() *Buffer
```
Get returns a buffer from the pool.

#### func (BufferFactory) MustPut

```go
func (cf BufferFactory) MustPut(buf *Buffer)
```
MustPut is the same as Put but panics if the buffer has already been returned to
the pool.

#### func (BufferFactory) Put

```go
func (cf BufferFactory) Put(buf *Buffer) error
```
Put returns the buffer to the pool. It returns an error if the buffer has
already been returned to the pool.

#### type BufferFactoryByteBufferCompat

```go
type BufferFactoryByteBufferCompat struct {
	BufferFactory
}
```


#### func (BufferFactoryByteBufferCompat) Get

```go
func (bf BufferFactoryByteBufferCompat) Get() ByteBuffer
```

#### func (BufferFactoryByteBufferCompat) Put

```go
func (bf BufferFactoryByteBufferCompat) Put(buf ByteBuffer)
```

#### type BufferFactoryInterfaceCompat

```go
type BufferFactoryInterfaceCompat struct {
	BufferFactory
}
```


#### func (BufferFactoryInterfaceCompat) Put

```go
func (b BufferFactoryInterfaceCompat) Put(buf *Buffer)
```

#### type ByteBuffer

```go
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
```

ByteBuffer is satisfied by [*pool.Buffer] and [*bytes.Buffer].

Note, we can't include Reset and Grow as our implementations return an error
while a [*bytes.Buffer] does not.

#### type Pool

```go
type Pool[T any] interface {
	Get() T
	Put(T)
}
```


#### type String

```go
type String struct {
	*strings.Builder
}
```


#### func (String) Cap

```go
func (s String) Cap() int
```

#### func (String) Grow

```go
func (s String) Grow(n int) error
```

#### func (String) Len

```go
func (s String) Len() int
```

#### func (String) MustGrow

```go
func (s String) MustGrow(n int)
```

#### func (String) MustLen

```go
func (s String) MustLen() int
```

#### func (String) MustReset

```go
func (s String) MustReset()
```

#### func (String) MustString

```go
func (s String) MustString() string
```

#### func (String) MustWriteString

```go
func (s String) MustWriteString(str string)
```
MustWriteString means Must Write String, like WriteString but will panic on
error.

#### func (String) Reset

```go
func (s String) Reset() error
```

#### func (String) String

```go
func (s String) String() string
```

#### func (String) Write

```go
func (s String) Write(p []byte) (int, error)
```

#### func (String) WriteByte

```go
func (s String) WriteByte(c byte) error
```

#### func (String) WriteRune

```go
func (s String) WriteRune(r rune) (int, error)
```

#### func (String) WriteString

```go
func (s String) WriteString(str string) (int, error)
```

#### type StringFactory

```go
type StringFactory struct {
}
```


#### func  NewStringFactory

```go
func NewStringFactory() StringFactory
```
NewStringFactory creates a new strings.Builder pool.

#### func (StringFactory) Get

```go
func (sf StringFactory) Get() *String
```
Get returns a strings.Builder from the pool.

#### func (StringFactory) MustPut

```go
func (sf StringFactory) MustPut(buf *String)
```

#### func (StringFactory) Put

```go
func (sf StringFactory) Put(buf *String) error
```
Put returns a strings.Builder back into to the pool after resetting it.

#### type WithPutError

```go
type WithPutError[T any] interface {
	Get() T
	Put(T) error
}
```

---
