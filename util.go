package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/rs/zerolog/log"
)

// Fprint is fmt.Fprint with error handling.
func Fprint(w io.Writer, s string) {
	_, err := fmt.Fprint(w, s)
	if err != nil {
		log.Error().Str("data", s).Err(err).Msg("Fprint failed!")
	}
}

// Fprintf is fmt.Fprintf with error handling.
func Fprintf(w io.Writer, format string, items ...any) {
	_, err := fmt.Fprintf(w, format, items...)
	if err != nil {
		log.Error().Str("fmt", format).Err(err).Msg("Fprintf failed!")
	}
}

// Abs will give you the positive version of a negative integer, quickly.
func Abs(n int) int {
	// ayyee smash 6ros
	n64 := int64(n)
	y := n64 >> 63
	n64 = (n64 ^ y) - y
	return int(n64)
}

// Float64ToBytes will take a float64 type and convert it to a slice of bytes.
func Float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// BytesToFloat64 will take a slice of bytes and convert it to a float64 type.
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
