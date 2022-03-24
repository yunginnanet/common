package hash

import (
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
)

// Blake2bSum ignores all errors and gives you a blakae2b 64 hash value as a byte slice. (or panics somehow)
func Blake2bSum(b []byte) []byte {
	Hasha, _ := blake2b.New(64, nil)
	Hasha.Write(b)
	return Hasha.Sum(nil)
}

// BlakeFileChecksum takes in the path of a file on the OS' filesystem and will attempt to calculate a blake2b checksum of the file's contents.
func BlakeFileChecksum(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	defer f.Close()

	empty := []byte{}

	if err != nil {
		return empty, err
	}

	buf, err := io.ReadAll(f)

	if err != nil {
		return empty, err
	}

	hasha, err := blake2b.New(64, nil)
	if err != nil {
		return empty, err
	}

	hasha.Write(buf)

	return hasha.Sum(nil), nil
}

// BlakeEqual will take in two byte slices, hash them with blake2b, and tell you if the resulting checksums match.
func BlakeEqual(a []byte, b []byte) bool {
	ahash := Blake2bSum(a)
	bhash := Blake2bSum(b)
	return string(ahash) == string(bhash)
}
