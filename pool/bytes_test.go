package pool

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewBufferFactory(t *testing.T) {
	t.Parallel()
	bf := NewBufferFactory()
	if bf.pool == nil {
		t.Fatalf("The pool is nil")
	}
}

func TestBufferFactory(t *testing.T) {
	t.Parallel()
	bf := NewBufferFactory()
	t.Run("BufferPut", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		if err := bf.Put(buf); err != nil {
			t.Fatalf("The buffer was not returned: %v", err)
		}
		if err := bf.Put(buf); err == nil {
			t.Fatalf("The buffer was returned twice")
		}
	})
	t.Run("BufferMustPut", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		bf.MustPut(buf)
		assertPanic(t, func() {
			bf.MustPut(buf)
		})
	})
	t.Run("BufferFactoryGet", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		if buf.Buffer == nil {
			t.Fatalf("The buffer is nil")
		}
		if buf.o == nil {
			t.Fatalf("The once is nil")
		}
	})
	t.Run("BufferBytes", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		if len(buf.Bytes()) != 0 {
			t.Fatalf("The bytes are not nil: %v", buf.Bytes())
		}
		buf.MustWrite([]byte("hello world"))
		if !bytes.Equal(buf.MustBytes(), []byte("hello world")) {
			t.Fatalf("The bytes are wrong")
		}
		bf.MustPut(buf)
		if buf.Bytes() != nil {
			t.Fatalf("The bytes are not nil")
		}
	})
	t.Run("BufferMustBytes", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_, err := buf.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf.MustBytes(), []byte("hello")) {
			t.Fatalf("The bytes are not equal")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustBytes()
		})
	})
	t.Run("BufferString", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		if buf.String() != "" {
			t.Fatalf("The string is not empty")
		}
		bf.MustPut(buf)
		if buf.String() != "" {
			t.Fatalf("The string is not empty")
		}
	})
	t.Run("BufferMustString", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_ = buf.MustString()
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustString()
		})
	})
	t.Run("BufferLen", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		if buf.Len() != 0 {
			t.Fatalf("The length is not zero")
		}
		bf.MustPut(buf)
		if buf.Len() != 0 {
			t.Fatalf("The length is not zero")
		}
	})
	t.Run("BufferMustLen", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_ = buf.MustLen()
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustLen()
		})
	})
	t.Run("BufferCap", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_ = buf.Cap()
		bf.MustPut(buf)
		if buf.Cap() != 0 {
			t.Fatalf("The capacity is not zero")
		}
	})
	t.Run("BufferReset", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		err := buf.Reset()
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 0 {
			t.Fatalf("The length is not zero")
		}
		bf.MustPut(buf)
	})
	t.Run("BufferMustReset", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		buf.MustReset()
		if buf.Len() != 0 {
			t.Fatalf("The length is not zero")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustReset()
		})
	})
	t.Run("BufferWrite", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_, err := buf.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 5 {
			t.Fatalf("The length is not five")
		}
		bf.MustPut(buf)
		written, werr := buf.Write([]byte("hello"))
		if written != 0 {
			t.Fatalf("The written is not zero")
		}
		if werr == nil {
			t.Fatalf("The error is nil")
		}
	})
	t.Run("BufferMustWrite", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		if buf.Len() != 5 {
			t.Fatalf("The length is not five")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustWrite([]byte("hello"))
		})
	})
	t.Run("BufferWriteByte", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		err := buf.WriteByte('h')
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 1 {
			t.Fatalf("The length is not one")
		}
		bf.MustPut(buf)
		werr := buf.WriteByte('h')
		if werr == nil {
			t.Fatalf("The error is nil")
		}
	})
	t.Run("BufferMustWriteByte", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWriteByte('h')
		if buf.Len() != 1 {
			t.Fatalf("The length is not one")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustWriteByte('h')
		})
	})
	t.Run("BufferWriteRune", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_, err := buf.WriteRune('h')
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 1 {
			t.Fatalf("The length is not one")
		}
		bf.MustPut(buf)
		written, werr := buf.WriteRune('h')
		if written != 0 {
			t.Fatalf("The written is not zero")
		}
		if werr == nil {
			t.Fatalf("The error is nil")
		}
	})
	t.Run("BufferMustWriteRune", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWriteRune('h')
		if buf.Len() != 1 {
			t.Fatalf("The length is not one")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustWriteRune('h')
		})
	})
	t.Run("BufferWriteString", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_, err := buf.WriteString("hello")
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 5 {
			t.Fatalf("The length is not five")
		}
		bf.MustPut(buf)
		written, werr := buf.WriteString("hello")
		if written != 0 {
			t.Fatalf("The written is not zero")
		}
		if werr == nil {
			t.Fatalf("The error is nil")
		}
	})
	t.Run("BufferGrow", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		err := buf.Grow(5)
		if buf.Cap() < 5 {
			t.Fatalf("The capacity is less than five: %d", buf.Cap())
		}
		if err != nil {
			t.Fatal(err)
		}
		bf.MustPut(buf)
		if buf.Cap() != 0 {
			t.Fatalf("The capacity is not zero")
		}
		if err = buf.Grow(1); err == nil {
			t.Fatal("The error is nil")
		}
	})
	t.Run("BufferTruncate", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		err := buf.Truncate(3)
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 3 {
			t.Fatalf("The length is not three")
		}
		if buf.String() != "hel" {
			t.Fatalf("The string is not hel")
		}
		bf.MustPut(buf)
	})
	t.Run("BufferMustTruncate", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		buf.MustTruncate(3)
		if buf.Len() != 3 {
			t.Fatalf("The length is not three")
		}
		if buf.String() != "hel" {
			t.Fatalf("The string is not hel")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustTruncate(3)
		})
	})
	t.Run("BufferRead", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		p := make([]byte, 5)
		n, err := buf.Read(p)
		if err != nil {
			t.Fatal(err)
		}
		if n != 5 {
			t.Fatalf("The n is not five")
		}
		if string(p) != "hello" {
			t.Fatalf("The string is not hello")
		}
		bf.MustPut(buf)
		if _, err = buf.Read(p); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferReadByte", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		b, err := buf.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if b != 'h' {
			t.Fatalf("The byte is not h")
		}
		bf.MustPut(buf)
		if _, err = buf.ReadByte(); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferReadRune", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		r, size, err := buf.ReadRune()
		if err != nil {
			t.Fatal(err)
		}
		if r != 'h' {
			t.Fatalf("The rune is not h")
		}
		if size != 1 {
			t.Fatalf("The size is not one")
		}
		bf.MustPut(buf)
		if _, _, err = buf.ReadRune(); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferUnreadByte", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		b, err := buf.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if b != 'h' {
			t.Fatalf("The byte is not h")
		}
		err = buf.UnreadByte()
		if err != nil {
			t.Fatal(err)
		}
		b, err = buf.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if b != 'h' {
			t.Fatalf("The byte is not h")
		}
		bf.MustPut(buf)
		if err = buf.UnreadByte(); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferUnreadRune", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		r, size, err := buf.ReadRune()
		if err != nil {
			t.Fatal(err)
		}
		if r != 'h' {
			t.Fatalf("The rune is not h")
		}
		if size != 1 {
			t.Fatalf("The size is not one")
		}
		err = buf.UnreadRune()
		if err != nil {
			t.Fatal(err)
		}
		r, size, err = buf.ReadRune()
		if err != nil {
			t.Fatal(err)
		}
		if r != 'h' {
			t.Fatalf("The rune is not h")
		}
		if size != 1 {
			t.Fatalf("The size is not one")
		}
		bf.MustPut(buf)
		if err = buf.UnreadRune(); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferReadBytes", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello!"))
		p, err := buf.ReadBytes('o')
		if err != nil {
			t.Fatal(err)
		}
		if string(p) != "hello" {
			t.Fatalf("The string is not hello: %v", string(p))
		}
		bf.MustPut(buf)
		if _, err = buf.ReadBytes('l'); err == nil {
			t.Fatal("The error is nil after returning the buffer")
		}
	})
	t.Run("BufferReadFrom", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		_, err := buf.ReadFrom(strings.NewReader("hello"))
		if err != nil {
			t.Fatal(err)
		}
		if buf.Len() != 5 {
			t.Fatalf("The length is not five")
		}
		bf.MustPut(buf)
		if _, err = buf.ReadFrom(strings.NewReader("hello")); err == nil {
			t.Fatal("The error is nil trying to use a returned buffer")
		}
		buf = bf.Get()
		buf.MustReadFrom(strings.NewReader("hello"))
		buf.MustTruncate(5)
		if buf.Len() != 5 {
			t.Fatalf("The length is not five")
		}
		bf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustReadFrom(strings.NewReader("hello"))
		})
	})
	t.Run("BufferWriteTo", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		n, err := buf.WriteTo(io.Discard)
		if err != nil {
			t.Fatal(err)
		}
		if n != 5 {
			t.Fatalf("The number of bytes is not five: %d", n)
		}
		bf.MustPut(buf)
		if _, err = buf.WriteTo(io.Discard); err == nil {
			t.Fatal("The error is nil trying to use a returned buffer")
		}
		assertPanic(t, func() {
			buf.MustWriteTo(io.Discard)
		})
	})
	t.Run("BufferNext", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		p := buf.Next(5)
		if string(p) != "hello" {
			t.Fatalf("The string is not hello")
		}
		bf.MustPut(buf)
		if p = buf.Next(5); p != nil {
			t.Fatalf("The slice is not nil")
		}
	})
	t.Run("NewSizedBufferFactory", func(t *testing.T) {
		t.Parallel()
		sized := NewSizedBufferFactory(4)
		buf := sized.Get()
		defer sized.MustPut(buf)
		if buf.Cap() != 4 {
			t.Errorf("Expected sized buffer from fresh factory to be cap == 4, got: %d", buf.Cap())
		}
	})
	t.Run("BufferWithParent", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		buf.MustWrite([]byte("world"))
		buf.MustWrite([]byte("!"))
		if buf.p != nil {
			t.Fatalf("The parent is not nil without assignment")
		}
		buf = buf.WithParent(&bf)
		if buf.p == nil {
			t.Errorf("The parent is nil after assignment")
		}
		if buf.p != &bf {
			t.Errorf("The parent is not the buffer factory")
		}
		if buf.String() != "helloworld!" {
			t.Fatalf("The string is not 'helloworld!': %v", buf.String())
		}
		bf.MustPut(buf)
	})
	t.Run("BufferClose", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf.MustWrite([]byte("hello"))
		buf.MustWrite([]byte("world"))
		buf.MustWrite([]byte("!"))
		if buf.IsClosed() {
			t.Fatalf("The buffer is closed before closing")
		}
		if err := buf.Close(); err == nil {
			t.Fatal("The error is nil after closing the buffer with no parent")
		}
		if buf.IsClosed() {
			t.Fatalf("The buffer is closed after failing to close")
		}
		if buf.String() != "helloworld!" {
			t.Fatalf("The string is not 'helloworld!' after unsuccessful close: %v", buf.String())
		}
		buf = buf.WithParent(&bf)
		if err := buf.Close(); err != nil {
			t.Fatal(err)
		}
		if !buf.IsClosed() {
			t.Fatalf("The buffer is not closed after closing")
		}
		if err := buf.Close(); err == nil {
			t.Fatal("The error is nil after closing an already closed buffer")
		}
	})
	t.Run("BufferCannotClose", func(t *testing.T) {
		t.Parallel()
		buf := bf.Get()
		buf = buf.WithParent(&bf)
		buf.MustWrite([]byte("hello"))
		buf.MustWrite([]byte("world"))
		buf.MustWrite([]byte("!"))
		if buf.String() != "helloworld!" {
			t.Fatalf("The string is not 'helloworld!': %v", buf.String())
		}
		bf.MustPut(buf)
		if buf.Close() == nil {
			t.Fatal("The error is nil after closing an already returned buffer with parent")
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("The panic is nil")
			}
		}()

		buf = nil
		if buf.Close() == nil {
			t.Fatal("The error is nil after closing a nil buffer")
		}
	})
	t.Run("BufferLockingSafety", func(t *testing.T) {
		// implement the Buffer.Safe* methods and test them
		t.Parallel()
		buf := bf.Get().WithMutex(&bf)
		if _, err := buf.SafeWrite([]byte("hello")); err != nil {
			t.Fatal(err)
		}
		if _, err := buf.SafeWrite([]byte("world")); err != nil {
			t.Fatal(err)
		}
		if _, err := buf.SafeWrite([]byte("!")); err != nil {
			t.Fatal(err)
		}
		if buf.SafeString() != "helloworld!" {
			t.Fatalf("The string is not 'helloworld!': %v", buf.String())
		}
		if b := buf.SafeNext(1); b == nil {
			t.Fatal("The byte is nil despite having written to the buffer")
		}
		if _, err := buf.SafeRead(make([]byte, 1)); err != nil {
			t.Fatal(err)
		}
		if err := buf.SafeUnreadByte(); err != nil {
			t.Fatal(err)
		}
		if _, err := buf.SafeWrite([]byte("yeet")); err != nil {
			t.Fatal(err)
		}
		go func() {
			if _, err := buf.SafeReadByte(); err != nil {
				t.Errorf("Error reading from buffer: %v", err)
			}
		}()
		for i := 0; i < 100; i++ {
			go func() {
				if _, err := buf.SafeWrite([]byte("yeet")); err != nil {
					t.Errorf("Error writing to buffer: %v", err)
				}
			}()
		}
		for i := 0; i < 100; i++ {
			go func() {
				// must be ran with race detector enabled or this check does nothing
				_ = buf.SafeLen()
			}()
		}
		for i := 0; i < 100; i++ {
			if _, err := buf.SafeRead([]byte("yeet")); err != nil {
				t.Errorf("Error reading from buffer: %v", err)
			}
		}
		for i := 0; i < 100; i++ {
			go func() {
				_ = buf.SafeBytes()
			}()
		}
		for i := 0; i < 100; i++ {
			go func() {
				_ = buf.SafeString()
			}()
		}

		for i := 0; i < 100; i++ {
			go func() {
				if _, err := buf.SafeWriteTo(io.Discard); err != nil {
					t.Errorf("Error writing to buffer: %v", err)
				}
			}()
		}
		for i := 0; i < 100; i++ {
			go func() {
				br := bytes.NewReader([]byte("yeet"))
				if _, err := buf.SafeReadFrom(br); err != nil {
					t.Errorf("Error reading from buffer: %v", err)
				}
			}()
		}
		for i := 0; i < 100; i++ {
			if _, err := buf.SafeReadFrom(strings.NewReader("yeet")); err != nil {
				t.Errorf("Error reading from buffer: %v", err)
			}
		}
		if err := buf.SafeReset(); err != nil {
			t.Errorf("Error reading from buffer: %v", err)
		}
		if _, err := buf.SafeWrite([]byte("hello")); err != nil {
			t.Fatal(err)
		}
		if err := buf.SafeClose(); err != nil {
			t.Fatal(err)
		}
		if err := buf.SafeClose(); err == nil {
			t.Fatal("The error is nil despite closing the buffer twice")
		}
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The error is nil despite closing the buffer twice")
				}
			}()
			_ = buf.SafeBytes()
		}()
		buf2 := bf.Get().WithMutex(&bf)
		_, err := buf2.SafeWrite([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
		b, err := buf2.SafeReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if b != 'h' {
			t.Fatalf("The byte is not 'h': %v", b)
		}
		bts := buf2.SafeNext(buf2.Len())
		if bts == nil {
			t.Fatal("The bytes are nil despite having written to the buffer")
		}
		if !bytes.Equal(bts, []byte("ello")) {
			t.Fatalf("The bytes are not 'ello': %v", bts)
		}
		buf2.SafeTruncate(0)
		if err = bf.SafePut(buf2); err != nil {
			t.Fatal(err)
		}
		if err = bf.SafePut(buf2); err == nil {
			t.Fatal("The error is nil despite putting the buffer twice")
		}
	})
	t.Run("BufferLockingSafetyMustFail", func(t *testing.T) {
		// implement the Buffer.Safe* methods and test them
		t.Parallel()
		buf := bf.Get()
		if _, err := buf.SafeWrite([]byte("hello")); err == nil {
			t.Fatal("the error is nil despite nil mutex")
		}
		if _, err := buf.SafeWrite([]byte("world")); err == nil {
			t.Fatal(err)
		}
		if _, err := buf.SafeWrite([]byte("!")); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			if buf.SafeString() != "helloworld!" {
				t.Fatalf("The string is not 'helloworld!': %v", buf.String())
			}
		}()
		if _, err := buf.SafeRead(make([]byte, 1)); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		if err := buf.SafeUnreadByte(); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		if _, err := buf.SafeWrite([]byte("yeet")); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		if _, err := buf.SafeReadByte(); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			_ = buf.SafeLen()
		}()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			_ = buf.SafeNext(1)
		}()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			_ = buf.SafeBytes()
		}()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			_ = buf.SafeString()
		}()
		if _, err := buf.SafeWriteTo(io.Discard); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		br := bytes.NewReader([]byte("yeet"))
		if _, err := buf.SafeReadFrom(br); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		if _, err := buf.SafeReadFrom(strings.NewReader("yeet")); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		if err := buf.SafeReset(); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			buf.SafeTruncate(1)
		}()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			_ = buf.SafeNext(1)
		}()
		if err := buf.SafeClose(); err == nil {
			t.Fatal("The error is nil despite nil mutex")
		}
		buf2 := bf.Get().WithMutex(&bf)
		buf2.Buffer = nil
		if err := buf2.SafeReset(); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if n := buf2.SafeLen(); n != 0 {
			t.Fatalf("The length is not 0: %v", n)
		}
		if _, err := buf2.SafeWriteTo(io.Discard); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if _, err := buf2.SafeReadFrom(br); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if _, err := buf2.SafeWrite([]byte("yeet")); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if _, err := buf2.SafeRead(make([]byte, 1)); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if err := buf2.SafeUnreadByte(); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if _, err := buf2.SafeReadByte(); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		if err := buf2.SafeClose(); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("The panic is nil")
				}
			}()
			buf2.SafeTruncate(1)
		}()
		if bts := buf2.SafeNext(1); bts != nil {
			t.Fatalf("The bytes are not nil: %v", bts)
		}
		buf2 = buf2.WithMutex(&bf)
		if bts := buf2.SafeBytes(); bts != nil {
			t.Fatalf("The bytes are not nil: %v", bts)
		}
		if str := buf2.SafeString(); str != "" {
			t.Fatalf("The string is not empty: %v", str)
		}
		if err := buf2.SafeClose(); err == nil {
			t.Fatal("The error is nil despite nil buffer")
		}
	})
}
