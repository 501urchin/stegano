// Package bits provides helpers for bitwise operations
package bits

import (
	"fmt"
)

// At gets the bit at b[idx] and returns it
func At(b byte, idx int) byte {
	if idx > 7 || idx < 0 {
		panic(fmt.Errorf("invalid bit index %d of 7", idx))
	}

	return (b >> idx) & 1
}

// Flip attempts to flip a bit at b[idx] and returns the new byte
func Flip(b byte, idx int) byte {
	if idx > 7 || idx < 0 {
		panic(fmt.Errorf("invalid bit index %d of 7", idx))
	}

	return b ^ (1 << idx)
}
