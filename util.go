package common

import (
	crip "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

// Fprint is fmt.Fprint with error handling.
func Fprint(w io.Writer, s string) {
	_, err := fmt.Fprint(w, s)
	if err != nil {
		log.Error().Str("data", s).Err(err).Msg("Fprint failed!")
	}
}

// RNG is a Random Number Generator
func RNG(n int) int {
	var seed int64
	err := binary.Read(crip.Reader, binary.BigEndian, &seed)
	if err != nil {
		panic(err)
	}
	rng := rand.New(rand.NewSource(seed))
	return rng.Intn(n)
}

// RandStr generates a random alphanumeric string with a max length of size.
func RandStr(size int) string {
	buf := make([]byte, size)
	for i := 0; i != size; i++ {
		buf[i] = charset[uint32(RNG(32))%uint32(len(charset))]
	}
	return string(buf)
}

// RandSleepMS sleeps for a random period of time with a maximum of n milliseconds.
func RandSleepMS(n int) {
	time.Sleep(time.Duration(RNG(n)) * time.Millisecond)
}

// Abs will give you the positive version of a negative integer, quickly.
func Abs(n int) int {
	// ayyee smash 6ros
	n64 := int64(n)
	y := n64 >> 63
	n64 = (n64 ^ y) - y
	return int(n64)
}

func OneInA(million int) bool {
	return RNG(million) == 1
}
