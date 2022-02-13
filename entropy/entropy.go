package entropy

import (
	crip "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

func RandomStrChoice(choice []string) string {
	strlen := len(choice)
	n := uint32(0)
	if strlen > 0 {
		n = uint32(RNG(16)) % uint32(strlen)
	}
	return choice[n]
}

func RNG(n int) int {
	var seed int64
	binary.Read(crip.Reader, binary.BigEndian, &seed)
	rng := rand.New(rand.NewSource(seed))
	return rng.Intn(n)
}

func OneInA(million int) bool {
	return RNG(million) == 1
}

// RandSleepMS sleeps for a random period of time with a maximum of n milliseconds.
func RandSleepMS(n int) {
	time.Sleep(time.Duration(RNG(n)) * time.Millisecond)
}

// characters used for the gerneration of random strings.
const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

// RandStr generates a random alphanumeric string with a max length of size. Charset used is all lowercase.
func RandStr(size int) string {
	buf := make([]byte, size)
	for i := 0; i != size; i++ {
		buf[i] = charset[uint32(RNG(32))%uint32(len(charset))]
	}
	return string(buf)
}
