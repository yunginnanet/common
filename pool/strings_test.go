package pool

import (
	"testing"
	"time"
)

func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestStringFactoryPanic(t *testing.T) {
	t.Parallel()
	sf := NewStringFactory()
	t.Run("StringsMustWrite", func(t *testing.T) {
		buf := sf.Get()
		buf.MustWriteString("hello world")
		if buf.Len() == 0 {
			t.Fatalf("The buffer is empty after we wrote to it")
		}
		if buf.String() != "hello world" {
			t.Fatalf("The buffer has the wrong content")
		}
	})
	t.Run("StringsMustWritePanic", func(t *testing.T) {
		t.Parallel()
		var badString *string = nil
		buf := sf.Get()
		assertPanic(t, func() {
			buf.MustWriteString(*badString)
		})
		/*		assertPanic(t, func() {
				buf.MustWriteString("")
			})*/
		if err := sf.Put(buf); err != nil {
			t.Fatalf("The buffer was not returned: %v", err)
		}
	})
	t.Run("StringsMustString", func(t *testing.T) {
		t.Parallel()
		buf := sf.Get()
		buf.MustWriteString("hello world")
		if buf.MustString() != "hello world" {
			t.Fatalf("The buffer has the wrong content")
		}
		sf.MustPut(buf)
		assertPanic(t, func() {
			buf.MustString()
		})
	})
	t.Run("StringsMust", func(t *testing.T) {
		t.Parallel()
		buf := sf.Get()
		buf.MustReset()
		_ = buf.MustLen()
		buf.MustGrow(10)
		err := sf.Put(buf)
		if err != nil {
			t.Fatalf("The buffer was not returned: %v", err)
		}
		assertPanic(t, func() {
			sf.MustPut(buf)
		})
		assertPanic(t, func() {
			buf.MustWriteString("hello")
		})
		assertPanic(t, func() {
			buf.MustGrow(10)
		})
		assertPanic(t, func() {
			buf.MustLen()
		})
		assertPanic(t, func() {
			buf.MustReset()
		})
	})
}

func TestStringFactory(t *testing.T) {
	t.Parallel()
	s := NewStringFactory()
	t.Run("StringPoolHelloWorld", func(t *testing.T) {
		t.Parallel()
		buf := s.Get()
		if _, err := buf.WriteString("hello"); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "hello" {
			t.Fatal("unexpected string")
		}
		if err := buf.WriteByte(' '); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "hello " {
			t.Fatalf("unexpected string: %s", buf.String())
		}
		if _, err := buf.WriteRune('w'); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "hello w" {
			t.Fatalf("unexpected string: %s", buf.String())
		}
		if _, err := buf.Write([]byte("orld")); err != nil {
			t.Fatal(err)
		}
		if err := buf.Grow(1); err != nil {
			t.Fatal(err)
		}
		if buf.Cap() == 0 {
			t.Fatal("expected capacity, got 0")
		}
		if err := buf.Reset(); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "" {
			t.Fatalf("unexpected string: %s", buf.String())
		}
		if err := s.Put(buf); err != nil {
			t.Fatal(err)
		}
		if err := s.Put(buf); err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("StringPoolCheckGetLength", func(t *testing.T) {
		t.Parallel()
		buf := s.Get()
		if buf.Len() > 0 {
			t.Fatalf("StringFactory.Put() did not reset the buffer")
		}
		if err := s.Put(buf); err != nil {
			t.Fatal(err)
		}
		if err := s.Put(buf); err == nil {
			t.Fatalf("StringFactory.Put() should have returned an error after already returning the buffer")
		}
	})
	t.Run("StringPoolGrowBuffer", func(t *testing.T) {
		t.Parallel()
		buf := s.Get()
		if err := buf.Grow(16); err != nil {
			t.Fatal(err)
		}
		if buf.Cap() < 16 {
			t.Fatalf("expected capacity of 1, got %d", buf.Cap())
		}
		if err := s.Put(buf); err != nil {
			t.Fatal(err)
		}
		if err := buf.Grow(10); err == nil {
			t.Fatalf("StringFactory.Grow() should not work after returning the buffer")
		}
		if buf.Cap() != 0 {
			t.Fatalf("StringFactory.Cap() should return 0 after returning the buffer")
		}
	})
	t.Run("StringPoolCleanBuffer", func(t *testing.T) {
		t.Parallel()
		time.Sleep(25 * time.Millisecond)
		got := s.Get()
		if got.String() != "" {
			t.Fatalf("should have gotten a clean buffer")
		}
		if err := s.Put(got); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("StringPoolWriteStringToReturnedBuffer", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		s.MustPut(got)
		if _, err := got.WriteString("a"); err == nil {
			t.Fatalf("should not be able to write to a returned buffer")
		}
	})
	t.Run("StringPoolWriteRuneToReturnedBuffer", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		s.MustPut(got)
		if _, err := got.WriteRune('a'); err == nil {
			t.Fatalf("should not be able to write to a returned buffer")
		}
	})
	t.Run("StringPoolWriteByteToReturnedBuffer", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		s.MustPut(got)
		if err := got.WriteByte('a'); err == nil {
			t.Fatalf("should not be able to write to a returned buffer")
		}
	})
	t.Run("StringPoolWriteToReturnedBuffer", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		s.MustPut(got)
		if _, err := got.Write([]byte("a")); err == nil {
			t.Fatalf("should not be able to write to a returned buffer")
		}
	})
	t.Run("StringPoolResetReturnedBuffer", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		s.MustPut(got)
		if err := got.Reset(); err == nil {
			t.Fatalf("should not be able to reset a returned buffer")
		}
		if str := got.String(); str != "" {
			t.Fatalf("should not be able to get string from a returned buffer")
		}
		if got.Len() != 0 {
			t.Fatalf("should not be able to write to a returned buffer")
		}
	})
	t.Run("StringFactoryMustNotPanicOnEmptyString", func(t *testing.T) {
		t.Parallel()
		got := s.Get()
		n, err := got.WriteString("")
		if err != nil {
			t.Fatal(err)
		}
		if n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
		if str := got.String(); str != "" {
			t.Fatalf("expected empty string, got %s", str)
		}
		if err := s.Put(got); err != nil {
			t.Fatal(err)
		}
		got = s.Get()
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unexpected panic: %v", r)
			}
		}()
		got.MustWriteString("")
		s.MustPut(got)
	})
}
