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
		n = getRandomUint32() % uint32(strlen)
	}
	return choice[n]
}

func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}

func RNG(n int) int {
	var seed int64
	binary.Read(crip.Reader, binary.BigEndian, &seed)
	rng := rand.New(rand.NewSource(seed))
	return rng.Intn(n)
}
