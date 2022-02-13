package entropy

import (
	"strings"
	"testing"
)

func Test_RNG(t *testing.T) {
	for n := 0; n != 500; n++ {
		zero := RNG(55555)
		one := RNG(55555)
		t.Logf("Random0: %d Random1: %d", zero, one)
		if zero == one {
			t.Errorf("RNG hit a duplicate! %d == %d", zero, one)
		}
		zero = 0
		one = 0
	}

}

func randStrChecks(zero, one string, t *testing.T) {
	if len(zero) != len(one) {
		t.Fatalf("RandStr output length inconsistency, len(zero) is %d but wanted len(one) which is %d", len(zero), len(one))
	}
	if len(zero) != 55 || len(one) != 55 {
		t.Fatalf("RandStr output length inconsistency, len(zero) is %d and len(one) is %d, but both should have been 55", len(zero), len(one))
	}
	if zero == one {
		t.Fatalf("RandStr hit a duplicate, %s == %s", zero, one)
	}
}

func Test_RandStr(t *testing.T) {
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		t.Logf("Random0: %s Random1: %s", zero, one)
		randStrChecks(zero, one, t)
		zero = ""
		one = ""
	}

}

func Test_RandStr_Entropy(t *testing.T) {
	for n := 0; n != 500; n++ {
		zero := RandStr(55)
		one := RandStr(55)
		randStrChecks(zero, one, t)
		zeroSplit := strings.Split(zero, "")
		oneSplit := strings.Split(one, "")
		var similarity = 0
		for i, char := range zeroSplit {
			if oneSplit[i] != char {
				continue
			}
			similarity++
			t.Logf("[-] zeroSplit[%d] is the same as oneSplit[%d] (%s)", i, i, char)
		}
		if similarity*4 > 55 {
			t.Errorf("[ENTROPY FAILURE] more than a quarter of the string is the same!\n zero: %s \n one: %s \nTotal similar: %d", zero, one, similarity)
		}
		t.Logf("[ENTROPY] Similarity score (lower is better): %d", similarity)
		zero = ""
		one = ""
	}

}
