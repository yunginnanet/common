package common

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
	"git.tcp.direct/kayos/common/hash"
	"git.tcp.direct/kayos/common/squish"
)

var needle = []byte(entropy.RandStr(16))

func TestBlakeEqualAndB64(t *testing.T) {
	var clone = make([]byte, len(needle))
	copy(clone, needle)

	if !hash.BlakeEqual(needle, clone) {
		t.Fatalf("BlakeEqual failed! Values %v and %v should have been equal.\n|---->Lengths: %d and %d",
			needle, clone, len(needle), len(clone),
		)
	}

	falseclone := []byte(entropy.RandStr(16))
	if hash.BlakeEqual(needle, falseclone) {
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
	bogusRd2, bogusWrt2 := io.Pipe()
	t.Logf("\n")
	go func() {
		Fprint(io.MultiWriter(bogusWrt, os.Stdout), fmt.Sprintf("[PASS] based[0] = %s", string(based[0])))
		Fprintf(io.MultiWriter(bogusWrt2, os.Stdout), "\n[PASS] based[1] = %s", string(based[0]))
	}()
	_ = bogusWrt.CloseWithError(io.ErrClosedPipe)
	_ = bogusWrt2.CloseWithError(io.ErrClosedPipe)
	_, err := bogusRd.Read([]byte{})
	if err == nil {
		t.Fatalf("should have been an error...")
	}
	_, err = bogusRd2.Read([]byte{})
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

func TestCruisinInMy64(t *testing.T) {
	data := 420.69
	databytes := Float64ToBytes(data)
	if len(databytes) < 1 {
		t.Fatalf("Float64ToBytes has returned a zero length value")
	}
	result := BytesToFloat64(databytes)
	if result != data {
		t.Fatalf("BytesToFloat64 failed! wanted %v and got %v", data, result)
	}
	t.Logf("original float64: %v -> Float64ToBytes %v -> BytesToFloat64 %v", data, databytes, result)
}

type phonyWriter struct{}

var o = &sync.Once{}
var fprintStatus bool
var pw = new(phonyWriter)

func (p2 phonyWriter) Write(p []byte) (int, error) {
	var err = errors.New("closed")
	fprintStatus = false
	o.Do(func() {
		err = nil
		fprintStatus = true
	})
	if err == nil {
		return len(p), err
	}
	return 0, err
}

func TestFprint(t *testing.T) {
	Fprint(pw, "asdf")
	if fprintStatus != true {
		t.Fatal("first Fprint test should have succeeded")
	}
	Fprint(pw, "asdf")
	if fprintStatus != false {
		t.Fatal("second Fprint test should not have succeeded")
	}
}
