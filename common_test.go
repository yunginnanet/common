package common

import (
	"errors"
	"sync"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
)

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
	var pw = new(phonyWriter)
	Fprint(pw, "asdf")
	if fprintStatus != true {
		t.Fatal("first Fprint test should have succeeded")
	}
	Fprint(pw, "asdf")
	if fprintStatus != false {
		t.Fatal("second Fprint test should not have succeeded")
	}
	pw = new(phonyWriter)
	fprintStatus = false
	o = &sync.Once{}
	Fprintf(pw, "%s", "asdf")
	if fprintStatus != true {
		t.Fatal("first Fprint test should have succeeded")
	}
	Fprintf(pw, "%s", "asdf")
	if fprintStatus != false {
		t.Fatal("second Fprint test should not have succeeded")
	}
}
