package hash

import (
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
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
	TypeCRC64ISO
	TypeCRC64ECMA
)

var typeToString = map[Type]string{
	TypeNull: "null", TypeBlake2b: "blake2b", TypeSHA1: "sha1",
	TypeSHA256: "sha256", TypeSHA512: "sha512",
	TypeMD5: "md5", TypeCRC32: "crc32",
	TypeCRC64ISO: "crc64-iso", TypeCRC64ECMA: "crc64-ecma",
}

var stringToType = map[string]Type{
	"null": TypeNull, "blake2b": TypeBlake2b, "sha1": TypeSHA1,
	"sha256": TypeSHA256, "sha512": TypeSHA512,
	"md5": TypeMD5, "crc32": TypeCRC32,
	"crc64-iso": TypeCRC64ISO, "crc64-ecma": TypeCRC64ECMA,
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
	crc64ISOPool = &sync.Pool{
		New: func() interface{} {
			// ISO and ECMA are pre-computed in the stdlib, so Make is just fetching them, not computing them.
			h := crc64.New(crc64.MakeTable(crc64.ISO))
			return h
		},
	}
	crc64ECMAPool = &sync.Pool{
		New: func() interface{} {
			h := crc64.New(crc64.MakeTable(crc64.ECMA))
			return h
		},
	}
	bufPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 40)
		},
	}
)

func getHasher(ht Type) hash.Hash {
	switch ht {
	case TypeBlake2b:
		return blake2bPool.Get().(hash.Hash)
	case TypeSHA1:
		return sha1Pool.Get().(hash.Hash)
	case TypeSHA256:
		return sha256Pool.Get().(hash.Hash)
	case TypeSHA512:
		return sha512Pool.Get().(hash.Hash)
	case TypeMD5:
		return md5Pool.Get().(hash.Hash)
	case TypeCRC32:
		return crc32Pool.Get().(hash.Hash)
	case TypeCRC64ISO:
		return crc64ISOPool.Get().(hash.Hash)
	case TypeCRC64ECMA:
		return crc64ECMAPool.Get().(hash.Hash)
	default:
		return nil
	}
}

func putHasher(ht Type, h hash.Hash) {
	h.Reset()
	switch ht {
	case TypeBlake2b:
		blake2bPool.Put(h)
	case TypeSHA1:
		sha1Pool.Put(h)
	case TypeSHA256:
		sha256Pool.Put(h)
	case TypeSHA512:
		sha512Pool.Put(h)
	case TypeMD5:
		md5Pool.Put(h)
	case TypeCRC32:
		crc32Pool.Put(h)
	case TypeCRC64ISO:
		crc64ISOPool.Put(h)
	case TypeCRC64ECMA:
		crc64ECMAPool.Put(h)
	default:
	}
}

func Sum(ht Type, b []byte) []byte {
	h := getHasher(ht)
	h.Write(b)
	b2 := bufPool.Get().([]byte)[0:0]
	if cap(b2) < h.Size() {
		b2 = append(b2, make([]byte, h.Size()-cap(b2))...)
	}
	sum := h.Sum(b2[:h.Size()])
	putHasher(ht, h)
	return sum
}

// RecycleChecksum will return the given byteslice previously allocated by [Sum] or [SumFile] to the buffer pool.
// This is useful to reduce memory allocations when you are done with the byte slice,
// but you MUST NOT reference the byte slice after calling this function.
//
// These bytes will be zeroed upon subsequent use, however;
// you should not rely on this behavior if the checksum data is sensitive.
// Either zero the data yourself or do not use this function.
func RecycleChecksum(sum []byte) {
	bufPool.Put(sum)
}

// SumHex will return the hex-encoded checksum string of the given byteslice.
//
// Note that this function makes a copy of the checksum data when it creates a hex-encoded string.
// Afterwards, the checksum byte slice is returned to the buffer pool.
//
// In short: do not use this function for sensitive data. Just use [Sum] and encode the checksum yourself.
// See notes on [RecycleChecksum] for more information.
func SumHex(ht Type, b []byte) string {
	s := Sum(ht, b)
	oldLen := len(s)
	needLen := hex.EncodedLen(len(s))
	s = hex.AppendEncode(s, s)
	str := string(s[oldLen : oldLen+needLen])
	bufPool.Put(s)
	return str
}

// SumFile will attempt to calculate a blake2b checksum of the given file path's contents.
// It will read the entire file into memory and return the checksum.
// If the file is empty, an error will be returned.
func SumFile(ht Type, path string) (buf []byte, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	if closeErr := f.Close(); closeErr != nil {
		err = fmt.Errorf("failed to close file during BlakeFileChecksum: %w", closeErr)
	}

	h := getHasher(ht)

	buf = bufPool.Get().([]byte)

	switch {
	case len(buf) == cap(buf):
	case cap(buf) >= 64 && cap(buf) > len(buf):
		buf = buf[:cap(buf)]
	case len(buf) < 64:
		buf = append(buf, make([]byte, 1024)...)
	default:
	}

	n, err := io.CopyBuffer(h, f, buf)
	if err != nil {
		putHasher(ht, h)
		bufPool.Put(buf)
		return nil, err
	}
	if n == 0 {
		putHasher(ht, h)
		bufPool.Put(buf)
		return nil, errors.New("file is empty")
	}
	h.Sum(buf[:0])
	putHasher(ht, h)
	return buf, nil
}
