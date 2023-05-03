package hash

import (
	"bytes"
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"

	"golang.org/x/crypto/blake2b"
)

type Type int8

const (
	TypeNull Type = iota
	TypeBlake2b
	TypeSHA1
	TypeSHA256
	TypeSHA512
	TypeMD5
)

var (
	sha1Pool = &sync.Pool{
		New: func() interface{} {
			return sha1.New() //nolint:gosec
		},
	}
	sha256Pool = &sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
	sha512Pool = &sync.Pool{
		New: func() interface{} {
			return sha512.New()
		},
	}
	md5Pool = &sync.Pool{
		New: func() interface{} {
			return md5.New() //nolint:gosec
		},
	}
	blake2bPool = &sync.Pool{
		New: func() interface{} {
			h, _ := blake2b.New(blake2b.Size, nil)
			return h
		},
	}
)

func Sum(ht Type, b []byte) []byte {
	var h hash.Hash
	switch ht {
	case TypeBlake2b:
		h = blake2bPool.Get().(hash.Hash)
		defer blake2bPool.Put(h)
	case TypeSHA1:
		h = sha1Pool.Get().(hash.Hash)
		defer sha1Pool.Put(h)
	case TypeSHA256:
		h = sha256Pool.Get().(hash.Hash)
		defer sha256Pool.Put(h)
	case TypeSHA512:
		h = sha512Pool.Get().(hash.Hash)
		defer sha512Pool.Put(h)
	case TypeMD5:
		h = md5Pool.Get().(hash.Hash)
		defer md5Pool.Put(h)
	default:
		return nil
	}
	h.Write(b)
	sum := h.Sum(nil)
	h.Reset()
	return sum
}

// Blake2bSum ignores all errors and gives you a blakae2b 64 hash value as a byte slice. (or panics somehow)
func Blake2bSum(b []byte) []byte {
	return Sum(TypeBlake2b, b)
}

// BlakeFileChecksum will attempt to calculate a blake2b checksum of the given file path's contents.
func BlakeFileChecksum(path string) (buf []byte, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := f.Close(); err != nil {
			err = fmt.Errorf("failed to close file during BlakeFileChecksum: %w", closeErr)
		}
	}()

	buf, _ = io.ReadAll(f)
	if len(buf) == 0 {
		return nil, errors.New("file is empty")
	}

	return Sum(TypeBlake2b, buf), nil
}

// BlakeEqual will take in two byte slices, hash them with blake2b, and tell you if the resulting checksums match.
func BlakeEqual(a []byte, b []byte) bool {
	return bytes.Equal(Blake2bSum(a), Blake2bSum(b))
}
