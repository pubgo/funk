package stack

import (
	"os"
	"strings"
)

// InsideTest returns true inside a Go test
func InsideTest() bool {
	return len(os.Args) > 1 && strings.HasSuffix(os.Args[0], ".test") &&
		strings.HasPrefix(os.Args[1], "-test.")
}
