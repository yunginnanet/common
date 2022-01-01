package common

import (
	"fmt"
	"io"

	crip "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/blake2b"
)

// Fprint is fmt.Fprint with error handling.
func Fprint(w io.Writer, s string) {
	_, err := fmt.Fprint(w, s)
	if err != nil {
		log.Error().Str("data", s).Err(err).Msg("Fprint failed!")
	}
}

func bytesToHash(b []byte) []byte {
	Hasha, _ := blake2b.New(64, nil)
	Hasha.Write(b)
	return Hasha.Sum(nil)
}

func CompareChecksums(a []byte, b []byte) bool {
	ahash := bytesToHash(a)
	bhash := bytesToHash(b)
	return string(ahash) == string(bhash)
}

func RNG(n int) int {
	var seed int64
	err := binary.Read(crip.Reader, binary.BigEndian, &seed)
	if err != nil {
		panic(err)
	}
	rng := rand.New(rand.NewSource(seed))
	return rng.Intn(n)
}

func SnoozeMS(n int) {
	time.Sleep(time.Duration(RNG(n)) * time.Millisecond)
}

func Abs(n int) int {
	// ayyee smash 6ros
	n64 := int64(n)
	y := n64 >> 63
	n64 = (n64 ^ y) - y
	return int(n64)
}
