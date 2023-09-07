package entropy

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

var dupCount = 0

func check[T comparable](t *testing.T, zero T, one T) {
	t.Helper()
	if zero == one {
		dupCount++
		t.Errorf("hit a duplicate! %v == %v", zero, one)
		t.Logf("duplicates so far: %d", dupCount)
	}
}

func Test_RNG(t *testing.T) {
	t.Parallel()
	// for coverage
	setSharedRand()
	RandSleepMS(5)
	hardLocc.Lock()
	sharedRand = nil
	getSharedRand = &sync.Once{}
	hardLocc.Unlock()
	//  - - - - - -
	if OneInA(1000000) {
		println(string([]byte{
			0x66, 0x75, 0x63, 0x6B, 0x68,
			0x6F, 0x6C, 0x65, 0x20, 0x6A,
			0x6F, 0x6E, 0x65, 0x73, 0x2E,
		}))
	}

	for n := 0; n != 55555; n++ {
		check(t, RNG(123454321), RNG(123454321))
		check(t, RNGUint32(), RNGUint32())
	}
}

func Test_OneInA(t *testing.T) {
	t.Parallel()
	for n := 0; n < 100; n++ {
		yes := ""
		if OneInA(1) {
			yes = "hello"
		}
		if yes != "hello" {
			t.Fatalf("OneInA failed to trigger when provided '1' as an argument")
		}
	}
}

func randStrChecks(t *testing.T, zero, one string, intendedLength int) {
	t.Helper()
	if len(zero) != len(one) {
		t.Fatalf("RandStr output length inconsistency, len(zero) is %d but wanted len(one) which is %d", len(zero), len(one))
	}
	if len(zero) != intendedLength || len(one) != intendedLength {
		t.Fatalf(
			"RandStr output length inconsistency, "+
				"len(zero) is %d and len(one) is %d, but both should have been 55", len(zero), len(one))
	}
	check(t, zero, one)
}

func Test_RandStr(t *testing.T) {
	t.Parallel()
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		t.Logf("Random0: %s Random1: %s", zero, one)
		randStrChecks(t, zero, one, 55)
	}
	t.Logf("[SUCCESS] RandStr had no collisions")
}

func Test_RandStrWithUpper(t *testing.T) {
	t.Parallel()
	for n := 0; n != 500; n++ {
		zero := RandStrWithUpper(15)
		one := RandStrWithUpper(15)
		t.Logf("Random0: %s Random1: %s", zero, one)
		randStrChecks(t, zero, one, 15)
	}
	t.Logf("[SUCCESS] RandStr had no collisions")
}

func Test_RandStr_Entropy(t *testing.T) {
	t.Parallel()
	var totalScore = 0
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		randStrChecks(t, zero, one, 55)
		zeroSplit := strings.Split(zero, "")
		oneSplit := strings.Split(one, "")
		var similarity = 0
		for i, char := range zeroSplit {
			if oneSplit[i] != char {
				continue
			}
			similarity++
		}
		if similarity*4 > 55 {
			t.Errorf("[ENTROPY FAILURE] more than a quarter of the string is the same!\n "+
				"zero: %s \n one: %s \nTotal similar: %d",
				zero, one, similarity)
		}
		// t.Logf("[ENTROPY] Similarity score (lower is better): %d", similarity)
		totalScore += similarity
	}
	t.Logf("[ENTROPY] final score (lower is better): %d (RandStr)", totalScore)
}

func Test_RandomStrChoice(t *testing.T) {
	t.Parallel()
	if RandomStrChoice([]string{}) != "" {
		t.Fatalf("RandomStrChoice returned a value when given an empty slice")
	}
	var slice []string
	for n := 0; n != 500; n++ {
		slice = append(slice, RandStr(555))
	}
	check(t, RandomStrChoice(slice), RandomStrChoice(slice))
}

func Test_RNGUint32(t *testing.T) {
	t.Parallel()
	// start globals fresh, just for coverage.
	setSharedRand()
	getSharedRand = &sync.Once{}
	RNGUint32()
}

func Benchmark_RandStr(b *testing.B) {
	toTest := []int{5, 25, 55, 500, 55555}
	for _, n := range toTest {
		for i := 1; i != 5; i++ {
			b.Run(fmt.Sprintf("%dchar/run%d", n, i), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for tn := 0; tn != b.N; tn++ {
					RandStr(n)
				}
			})
		}
	}
}
