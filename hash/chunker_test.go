package hash

import (
	"crypto/sha256"
	"errors"
	"io"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestChunkedHasherNewWithNegativeChunkSize(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, -1)
	if h == nil {
		t.Errorf("NewChunkedHasher() expected a ChunkedHasher, got nil")
	}
}

func TestChunkedHasherWriteAfterClose(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	h.Close()
	_, err := h.Write([]byte("test"))
	if err != ErrHasherClosed {
		t.Errorf("Write() expected error %v, got %v", ErrHasherClosed, err)
	}
}

func TestChunkedHasherNextWithNoData(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	dst := make([]byte, 0, h.hashedChunkSize)
	_, _, err := h.Next(dst)
	if !errors.Is(err, io.EOF) {
		spew.Dump(h)
		t.Errorf("Next() expected eof, got %v", err)
	}
}

func TestChunkedHasherNextAfterClose(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	h.Close()
	dst := make([]byte, 0, h.hashedChunkSize)
	_, _, err := h.Next(dst)
	if err != ErrHasherClosed {
		t.Errorf("Next() expected error %v, got %v", ErrHasherClosed, err)
	}
}

func TestChunkedHasherFlushAfterClose(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	h.Close()
	err := h.Flush()
	if err != ErrHasherClosed {
		t.Errorf("Flush() expected error %v, got %v", ErrHasherClosed, err)
	}
}

func TestChunkedHasherFlushWithNoData(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	err := h.Flush()
	if err != nil {
		t.Errorf("Flush() error = %v", err)
	}
}

func TestChunkedHasherIngressChunkSize(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	if size := h.IngressChunkSize(); size != 10 {
		t.Errorf("IngressChunkSize() expected %v, got %v", 10, size)
	}
}

func TestChunkedHasherEgressChunkSize(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	if size := h.EgressChunkSize(); size != sha256.Size {
		t.Errorf("EgressChunkSize() expected %v, got %v", sha256.Size, size)
	}
}
