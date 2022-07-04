package hash

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"os"

	"github.com/pkg/errors"
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

func Sum(ht Type, b []byte) []byte {
	var h hash.Hash
	switch ht {
	case TypeBlake2b:
		h, _ := blake2b.New(blake2b.Size, nil)
		h.Write(b)
		return h.Sum(nil)
	case TypeSHA1:
		h = sha1.New()
	case TypeSHA256:
		h = sha256.New()
	case TypeSHA512:
		h = sha512.New()
	case TypeMD5:
		h = md5.New()
	default:
		return nil
	}
	h.Write(b)
	return h.Sum(nil)
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
			err = errors.Wrapf(err, "failed to close file: %s", closeErr)
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
