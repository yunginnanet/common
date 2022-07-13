package entropy

import (
	crip "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"

	"nullprogram.com/x/rng"
)

var (
	sharedRand    *rand.Rand
	getSharedRand = &sync.Once{}
)

// RandomStrChoice returns a random item from an input slice of strings.
func RandomStrChoice(choice []string) string {
	if len(choice) > 0 {
		return choice[RNGUint32()%uint32(len(choice))]
	}
	return ""
}

// GetCryptoSeed returns a random int64 derived from crypto/rand.
// This can be used as a seed for the math/rand package.
func GetCryptoSeed() int64 {
	var seed int64
	_ = binary.Read(crip.Reader, binary.BigEndian, &seed)
	return seed
}

// GetOptimizedRand returns a pointer to a *new* rand.Rand which uses crypto/rand to seed a splitmix64 rng.
// Does not use the global/shared instance of a splitmix64 rng, but instead creates a new one.
func GetOptimizedRand() *rand.Rand {
	r := new(rng.SplitMix64)
	r.Seed(GetCryptoSeed())
	return rand.New(r)
}

// RNGUint32 returns a random uint32 using crypto/rand and splitmix64.
func RNGUint32() uint32 {
	getSharedRand.Do(func() {
		sharedRand = GetOptimizedRand()
	})
	return sharedRand.Uint32()
}

/*RNG returns integer with a maximum amount of 'n' using a global/shared instance of a splitmix64 rng.
  - Benchmark_FastRandStr5-24            25205089      47.03 ns/op
  - Benchmark_FastRandStr25-24       	7113620     169.8  ns/op
  - Benchmark_FastRandStr55-24       	3520297     340.7  ns/op
  - Benchmark_FastRandStr500-24      	 414966    2837    ns/op
  - Benchmark_FastRandStr55555-24    	   3717  315229    ns/op

*/
func RNG(n int) int {
	getSharedRand.Do(func() {
		sharedRand = GetOptimizedRand()
	})
	return sharedRand.Intn(n)
}

// OneInA generates a random number with a maximum of 'million' (input int).
// If the resulting random number is equal to 1, then the result is true.
func OneInA(million int) bool {
	return RNG(million) == 1
}

// RandSleepMS sleeps for a random period of time with a maximum of n milliseconds.
func RandSleepMS(n int) {
	time.Sleep(time.Duration(RNG(n)) * time.Millisecond)
}

// characters used for the gerneration of random strings.
const charset = "abcdefghijklmnopqrstuvwxyz1234567890"
const charsetWithUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

// RandStr generates a random alphanumeric string with a max length of size.
// Alpha charset used is a-z all lowercase.
func RandStr(size int) string {
	buf := make([]byte, size)
	for i := 0; i != size; i++ {
		buf[i] = charset[uint32(RNG(36))%uint32(len(charset))]
	}
	return string(buf)
}

// RandStrWithUpper generates a random alphanumeric string with a max length of size.
// Alpha charset used is a-Z mixed case.
func RandStrWithUpper(size int) string {
	buf := make([]byte, size)
	for i := 0; i != size; i++ {
		buf[i] = charsetWithUpper[uint32(RNG(62))%uint32(len(charsetWithUpper))]
	}
	return string(buf)
}
