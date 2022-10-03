package pool

import (
	"testing"
	"time"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestStringFactoryPanic(t *testing.T) {
	s := NewStringFactory()
	t.Run("StringsMustWrite", func(t *testing.T) {
		buf := s.Get()
		buf.MustWstr("hello world")
		if buf.Len() == 0 {
			t.Fatalf("The buffer is empty after we wrote to it")
		}
		if buf.String() != "hello world" {
			t.Fatalf("The buffer has the wrong content")
		}
	})
	t.Run("StringsMustWritePanic", func(t *testing.T) {
		var badString *string = nil
		buf := s.Get()
		assertPanic(t, func() {
			buf.MustWstr(*badString)
		})
		assertPanic(t, func() {
			buf.MustWstr("")
		})
		if err := s.Put(buf); err != nil {
			t.Fatalf("The buffer was not returned: %v", err)
		}
	})
	t.Run("StringsPanic", func(t *testing.T) {
		buf := s.Get()
		err := s.Put(buf)
		if err != nil {
			t.Fatalf("The buffer was not returned: %v", err)
		}
		assertPanic(t, func() {
			s.MustPut(buf)
		})
		assertPanic(t, func() {
			buf.MustWstr("hello")
		})
	})
}

func TestStringFactory(t *testing.T) {
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
		if err := buf.Grow(1); err != nil {
			t.Fatal(err)
		}
		if buf.Cap() != 1 {
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
		time.Sleep(100 * time.Millisecond)
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
}
