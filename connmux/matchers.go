package connmux

import (
	"bytes"
	"io"
)

// Any matches any connection.
func Any() Matcher {
	return func(io.Reader) bool { return true }
}

// Prefix matches if the connection starts with any of the provided prefixes.
//
// This matcher only reads as many bytes as the longest prefix.
func Prefix(prefixes ...[]byte) Matcher {
	// Copy to avoid surprising caller mutations.
	ps := make([][]byte, 0, len(prefixes))
	m := 0
	for _, p := range prefixes {
		if len(p) == 0 {
			continue
		}
		cp := append([]byte(nil), p...)
		ps = append(ps, cp)
		if len(cp) > m {
			m = len(cp)
		}
	}
	return func(r io.Reader) bool {
		if m == 0 {
			return false
		}
		buf := make([]byte, m)
		n, err := io.ReadFull(r, buf)
		if err != nil {
			return false
		}
		buf = buf[:n]
		for _, p := range ps {
			if len(buf) >= len(p) && bytes.Equal(buf[:len(p)], p) {
				return true
			}
		}
		return false
	}
}
