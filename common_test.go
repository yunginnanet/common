package common

import (
	"fmt"
	"io"
	"os"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
	"git.tcp.direct/kayos/common/hash"
	"git.tcp.direct/kayos/common/squish"
)

var needle = []byte(entropy.RandStr(16))

func TestBlakeEqualAndB64(t *testing.T) {
	var clone = make([]byte, len(needle))
	for i, c := range needle {
		clone[i] = c
	}
	if !hash.BlakeEqual(needle, clone) {
		t.Fatalf("BlakeEqual failed! Values %v and %v should have been equal.\n|---->Lengths: %d and %d",
			needle, clone, len(needle), len(clone),
		)
	}
	clone = make([]byte, len(needle))
	clone = []byte(entropy.RandStr(16))
	if hash.BlakeEqual(needle, clone) {
		t.Fatalf("BlakeEqual failed! Values %v and %v should NOT have been equal.\n|---->Lengths: %d and %d",
			needle, clone, len(needle), len(clone),
		)
	}

	var based = [2][]byte{needle, clone}

	based[0] = []byte(squish.B64e(based[0]))
	based[1] = []byte(squish.B64e(based[0]))

	if hash.BlakeEqual(based[0], based[1]) {
		t.Fatalf("Base64 encoding failed! Values %v and %v should NOT have been equal.\n|---->Lengths: %d and %d",
			based[0], based[1], len(based[0]), len(based[1]),
		)
	}

	// sneakin in some code coverage rq dwai nbd
	bogusRd, bogusWrt := io.Pipe()
	t.Logf("\n")
	go func() {
		Fprint(io.MultiWriter(bogusWrt, os.Stdout), fmt.Sprintf("[PASS] based[0] = %s\n[PASS] based[1] = %s", string(based[0]), string(based[1])))
	}()
	_ = bogusWrt.CloseWithError(io.ErrClosedPipe)
	_, err := bogusRd.Read([]byte{})
	if err == nil {
		t.Fatalf("should have been an error...")
	}
}

func TestAbs(t *testing.T) {
	var start = int32(entropy.RNG(5))
	for start < 1 {
		t.Logf("Re-rolling for a non-zero value... %d", start)
		start = int32(entropy.RNG(5))
	}

	less := start * 2
	negged := start - less

	if negged == start {
		t.Fatalf("the sky is falling. up is down: %d should not equal %d.", start, negged)
	}

	if Abs(int(negged)) != int(start) {
		t.Fatalf("Abs failed! values %d and %d should have been equal.", start, negged)
	}
}
