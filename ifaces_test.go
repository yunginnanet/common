package common

import (
	"net"
	"testing"
)

func needsDialer(t *testing.T, d any) {
	if _, ok := d.(Dialer); !ok {
		t.Fatal("d is not a Dialer")
	}
}

func TestDialer(t *testing.T) {
	needsDialer(t, &net.Dialer{})
}
