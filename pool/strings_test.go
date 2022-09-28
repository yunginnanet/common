package pool

import (
	"testing"
)

func TestStringFactory(t *testing.T) {
	s := NewStringFactory()
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
	if s.Get().Len() > 0 {
		t.Fatalf("StringFactory.Put() did not reset the buffer")
	}
	if err := s.Put(buf); err == nil {
		t.Fatalf("StringFactory.Put() should have returned an error after already returning the buffer")
	}
	if err := buf.Grow(10); err == nil {
		t.Fatalf("StringFactory.Grow() should not work after returning the buffer")
	}
	if buf.Cap() != 0 {
		t.Fatalf("StringFactory.Cap() should return 0 after returning the buffer")
	}
	got := s.Get()
	if got.String() != "" {
		t.Fatalf("should have gotten a clean buffer")
	}
	if err := s.Put(got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := got.WriteString("a"); err == nil {
		t.Fatalf("should not be able to write to a returned buffer")
	}
	if _, err := got.WriteRune('a'); err == nil {
		t.Fatalf("should not be able to write to a returned buffer")
	}
	if err := got.WriteByte('a'); err == nil {
		t.Fatalf("should not be able to write to a returned buffer")
	}
	if _, err := got.Write([]byte("a")); err == nil {
		t.Fatalf("should not be able to write to a returned buffer")
	}
	if err := got.Reset(); err == nil {
		t.Fatalf("should not be able to reset a returned buffer")
	}
	if str := got.String(); str != "" {
		t.Fatalf("should not be able to get string from a returned buffer")
	}
	if got.Len() != 0 {
		t.Fatalf("should not be able to write to a returned buffer")
	}
	if got = s.Get(); got.Len() > 0 {
		t.Fatalf("StringFactory.Put() did not reset the buffer")
	}
	if got.Cap() != 0 {
		t.Fatalf("StringFactory.Put() did not reset the buffer")
	}
}
