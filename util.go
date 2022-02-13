package common

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

// Fprint is fmt.Fprint with error handling.
func Fprint(w io.Writer, s string) {
	_, err := fmt.Fprint(w, s)
	if err != nil {
		log.Error().Str("data", s).Err(err).Msg("Fprint failed!")
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
