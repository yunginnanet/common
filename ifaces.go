package common

import "net"

// Dialer is an interface that should exist in stdlib honestly. Make it make sense that it doesn't.
type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}
