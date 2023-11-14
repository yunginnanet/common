package hash

import (
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
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
	TypeCRC32
)

var typeToString = map[Type]string{
	TypeNull: "null", TypeBlake2b: "blake2b", TypeSHA1: "sha1",
	TypeSHA256: "sha256", TypeSHA512: "sha512",
	TypeMD5: "md5", TypeCRC32: "crc32",
}

var stringToType = map[string]Type{
	"null": TypeNull, "blake2b": TypeBlake2b, "sha1": TypeSHA1,
	"sha256": TypeSHA256, "sha512": TypeSHA512,
	"md5": TypeMD5, "crc32": TypeCRC32,
}

func StringToType(s string) Type {
	t, ok := stringToType[s]
	if !ok {
		return TypeNull
	}
	return t
}

func (t Type) String() string {
	s, ok := typeToString[t]
	if !ok {
		return "unknown"
	}
	return s
}

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
	crc32Pool = &sync.Pool{
		New: func() interface{} {
			return crc32.NewIEEE()
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
	case TypeCRC32:
		h = crc32Pool.Get().(hash.Hash)
		defer crc32Pool.Put(h)
	default:
		return nil
	}
	h.Write(b)
	sum := h.Sum(nil)
	h.Reset()
	return sum
}

// SumFile will attempt to calculate a blake2b checksum of the given file path's contents.
func SumFile(ht Type, path string) (buf []byte, err error) {
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

	return Sum(ht, buf), nil
}
