package hash

import "golang.org/x/crypto/blake2b"

// Blake2bSum ignores all errors and gives you a blakae2b 64 hash value as a byte slice. (or panics somehow)
func Blake2bSum(b []byte) []byte {
	Hasha, _ := blake2b.New(64, nil)
	Hasha.Write(b)
	return Hasha.Sum(nil)
}

// BlakeEqual will take in two byte slices, hash them with blake2b, and tell you if the resulting checksums match.
func BlakeEqual(a []byte, b []byte) bool {
	ahash := Blake2bSum(a)
	bhash := Blake2bSum(b)
	return string(ahash) == string(bhash)
}
