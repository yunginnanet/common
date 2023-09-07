package entropy

import (
	crip "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"

	"nullprogram.com/x/rng"

	"git.tcp.direct/kayos/common/pool"
)

type randPool struct {
	sync.Pool
}

func (p *randPool) Get() *rand.Rand {
	return p.Pool.Get().(*rand.Rand)
}

func (p *randPool) Put(r *rand.Rand) {
	p.Pool.Put(r)
}

var (
	lolXD = randPool{
		Pool: sync.Pool{
			New: func() interface{} {
				sm64 := new(rng.SplitMix64)
				sm64.Seed(GetCryptoSeed())
				prng := rand.New(sm64) //nolint:gosec
				return prng
			},
		},
	}
	hardLocc      = &sync.Mutex{}
	sharedRand    *rand.Rand
	getSharedRand = &sync.Once{}
)

func setSharedRand() {
	hardLocc.Lock()
	sharedRand = lolXD.Get()
	hardLocc.Unlock()
}

func AcquireRand() *rand.Rand {
	return lolXD.Get()
}

func ReleaseRand(r *rand.Rand) {
	lolXD.Put(r)
}

// RandomStrChoice returns a random item from an input slice of strings.
func RandomStrChoice(choice []string) string {
	if len(choice) > 0 {
		return choice[RNGUint32()%uint32(len(choice))]
	}
	return ""
}

// GetCryptoSeed returns a random int64 derived from crypto/rand.
// This can be used as a seed for various PRNGs.
func GetCryptoSeed() int64 {
	var seed int64
	_ = binary.Read(crip.Reader, binary.BigEndian, &seed)
	return seed
}

// GetOptimizedRand returns a pointer to a *new* rand.Rand which uses GetCryptoSeed to seed an rng.SplitMix64.
// Does not use the global/shared instance of a splitmix64 rng, but instead creates a new one.
func GetOptimizedRand() *rand.Rand {
	r := new(rng.SplitMix64)
	r.Seed(GetCryptoSeed())
	return rand.New(r) //nolint:gosec
}

// GetSharedRand returns a pointer to our shared optimized rand.Rand which uses crypto/rand to seed a splitmix64 rng.
// WARNING - RACY - This is not thread safe, and should only be used in a single-threaded context.
func GetSharedRand() *rand.Rand {
	getSharedRand.Do(func() {
		setSharedRand()
	})
	return sharedRand
}

// RNGUint32 returns a random uint32 using crypto/rand and splitmix64.
func RNGUint32() uint32 {
	r := lolXD.Get()
	ui := r.Uint32()
	lolXD.Put(r)
	return ui
}

/*
RNG returns integer with a maximum amount of 'n' using a global/shared instance of a splitmix64 rng.
  - Benchmark_FastRandStr5-24          25205089      47.03 ns/op
  - Benchmark_FastRandStr25-24       	7113620     169.8  ns/op
  - Benchmark_FastRandStr55-24       	3520297     340.7  ns/op
  - Benchmark_FastRandStr500-24      	 414966    2837    ns/op
  - Benchmark_FastRandStr55555-24    	   3717  315229    ns/op
*/
func RNG(n int) int {
	r := lolXD.Get()
	i := r.Intn(n)
	lolXD.Put(r)
	return i
}

// OneInA generates a random number with a maximum of 'million' (input int).
// If the resulting random number is equal to 1, then the result is true.
func OneInA(million int) bool {
	if million == 1 {
		return true
	}
	return RNG(million) == 1
}

// RandSleepMS sleeps for a random period of time with a maximum of n milliseconds.
func RandSleepMS(n int) {
	time.Sleep(time.Duration(RNG(n)) * time.Millisecond)
}

// characters used for the gerneration of random strings.
const charset = "abcdefghijklmnopqrstuvwxyz1234567890"
const charsetWithUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

var strBufs = pool.NewStringFactory()

// RandStr generates a random alphanumeric string with a max length of size.
// Alpha charset used is a-z all lowercase.
func RandStr(size int) string {
	return randStr(charset, size)
}

// RandStrWithUpper generates a random alphanumeric string with a max length of size.
// Alpha charset used is a-Z mixed case.
func RandStrWithUpper(size int) string {
	return randStr(charsetWithUpper, size)
}

func randStr(chars string, size int) string {
	buf := strBufs.Get()
	r := lolXD.Get()
	for i := 0; i != size; i++ {
		ui32 := int(r.Uint32())
		_, _ = buf.WriteRune(rune(chars[ui32%len(chars)]))
	}
	lolXD.Put(r)
	s := buf.String()
	strBufs.MustPut(buf)
	return s
}
