package hash

import (
	"bytes"
	"crypto/sha256"
	"sync"
	"testing"
)

func TestChunkedHasher(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		chunkSize int
	}{
		{
			name:      "Exact chunk size",
			data:      []byte("1234567890"),
			chunkSize: 10,
		},
		{
			name:      "Less than chunk size",
			data:      []byte("12345"),
			chunkSize: 10,
		},
		{
			name:      "More than chunk size",
			data:      []byte("123456789012345"),
			chunkSize: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewChunkedHasher(TypeSHA256, tt.chunkSize)
			_, err := h.Write(tt.data)
			if err != nil {
				t.Errorf("Write() error = %v", err)
			}
			if err := h.Flush(); err != nil {
				t.Errorf("Flush() error = %v", err)
			}
			dst := make([]byte, 0, tt.chunkSize)
			dst, n, err := h.Next(dst)
			if err != nil {
				t.Errorf("Next() error = %v", err)
			}
			if n != tt.chunkSize {
				t.Errorf("Next() expected %d bytes, got %d", tt.chunkSize, n)
			}
			expectedHash := sha256.Sum256(tt.data)
			if !bytes.Equal(dst[:n], expectedHash[:]) {
				t.Errorf("Next() expected hash %x, got %x", expectedHash[:], dst[:n])
			}
		})
	}
}

func TestChunkedHasher_Closed(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	h.Close()
	_, err := h.Write([]byte("test"))
	if err != ErrHasherClosed {
		t.Errorf("Write() expected error %v, got %v", ErrHasherClosed, err)
	}
	_, _, err = h.Next(make([]byte, 10))
	if err != ErrHasherClosed {
		t.Errorf("Next() expected error %v, got %v", ErrHasherClosed, err)
	}
}

func TestChunkedHasher_Concurrent(t *testing.T) {
	h := NewChunkedHasher(TypeSHA256, 10)
	data := []byte("123456789012345")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := h.Write(data)
			if err != nil {
				t.Errorf("Write() error = %v", err)
			}
		}()
	}
	wg.Wait()
	if err := h.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}
	dst := make([]byte, 0, 10)
	dst, n, err := h.Next(dst)
	if err != nil {
		t.Errorf("Next() error = %v", err)
	}
	if n != 10 {
		t.Errorf("Next() expected %d bytes, got %d", 10, n)
	}
	expectedHash := sha256.Sum256(data)
	if !bytes.Equal(dst[:n], expectedHash[:]) {
		t.Errorf("Next() expected hash %x, got %x", expectedHash[:], dst[:n])
	}
}
