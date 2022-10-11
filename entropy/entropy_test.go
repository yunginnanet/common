package entropy

import (
	"strings"
	"sync"
	"testing"
)

func check[T comparable](zero T, one T, t *testing.T) {
	if zero == one {
		t.Errorf("hit a duplicate! %v == %v", zero, one)
	}
}

func Test_RNG(t *testing.T) {
	// for coverage
	sharedRand = GetSharedRand()
	RandSleepMS(5)
	sharedRand = nil
	getSharedRand = &sync.Once{}
	//  - - - - - -
	if OneInA(1000000) {
		println(string([]byte{
			0x66, 0x75, 0x63, 0x6B, 0x68,
			0x6F, 0x6C, 0x65, 0x20, 0x6A,
			0x6F, 0x6E, 0x65, 0x73, 0x2E,
		}))
	}

	for n := 0; n != 500; n++ {
		check(RNG(55555), RNG(55555), t)
		check(RNGUint32(), RNGUint32(), t)
	}
}

func randStrChecks(zero, one string, t *testing.T, intendedLength int) {
	if len(zero) != len(one) {
		t.Fatalf("RandStr output length inconsistency, len(zero) is %d but wanted len(one) which is %d", len(zero), len(one))
	}
	if len(zero) != intendedLength || len(one) != intendedLength {
		t.Fatalf("RandStr output length inconsistency, len(zero) is %d and len(one) is %d, but both should have been 55", len(zero), len(one))
	}
	check(zero, one, t)
}

func Test_RandStr(t *testing.T) {
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		t.Logf("Random0: %s Random1: %s", zero, one)
		randStrChecks(zero, one, t, 55)
	}
	t.Logf("[SUCCESS] RandStr had no collisions")
}

func Test_RandStrWithUpper(t *testing.T) {
	for n := 0; n != 500; n++ {
		zero := RandStrWithUpper(15)
		one := RandStrWithUpper(15)
		t.Logf("Random0: %s Random1: %s", zero, one)
		randStrChecks(zero, one, t, 15)
	}
	t.Logf("[SUCCESS] RandStr had no collisions")
}

func Test_RandStr_Entropy(t *testing.T) {
	var totalScore = 0
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		randStrChecks(zero, one, t, 55)
		zeroSplit := strings.Split(zero, "")
		oneSplit := strings.Split(one, "")
		var similarity = 0
		for i, char := range zeroSplit {
			if oneSplit[i] != char {
				continue
			}
			similarity++
			// t.Logf("[-] zeroSplit[%d] is the same as oneSplit[%d] (%s)", i, i, char)
		}
		if similarity*4 > 55 {
			t.Errorf("[ENTROPY FAILURE] more than a quarter of the string is the same!\n zero: %s \n one: %s \nTotal similar: %d",
				zero, one, similarity)
		}
		// t.Logf("[ENTROPY] Similarity score (lower is better): %d", similarity)
		totalScore += similarity
	}
	t.Logf("[ENTROPY] final score (lower is better): %d (RandStr)", totalScore)
}

func Test_RandomStrChoice(t *testing.T) {
	if RandomStrChoice([]string{}) != "" {
		t.Fatalf("RandomStrChoice returned a value when given an empty slice")
	}
	var slice []string
	for n := 0; n != 500; n++ {
		slice = append(slice, RandStr(555))
	}
	check(RandomStrChoice(slice), RandomStrChoice(slice), t)
}

func Test_RNGUint32(t *testing.T) {
	// start globals fresh, just for coverage.
	sharedRand = GetOptimizedRand()
	getSharedRand = &sync.Once{}
	RNGUint32()
}

func Benchmark_RandStr5(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStr(5)
	}
}

func Benchmark_RandStr25(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStr(25)
	}
}

func Benchmark_RandStr55(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStr(55)
	}
}

func Benchmark_RandStr500(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStr(500)
	}
}

func Benchmark_RandStr55555(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStr(55555)
	}
}

func Benchmark_RandStrWithUpper5(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStrWithUpper(5)
	}
}

func Benchmark_RandStrWithUpper25(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStrWithUpper(25)
	}
}

func Benchmark_RandStrWithUpper55(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStrWithUpper(55)
	}
}

func Benchmark_RandStrWithUpper500(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStrWithUpper(500)
	}
}

func Benchmark_RandStrWithUpper55555(b *testing.B) {
	for n := 0; n != b.N; n++ {
		RandStrWithUpper(55555)
	}
}
