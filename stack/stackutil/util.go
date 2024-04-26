package stackutil

import (
	"os"
	"strings"

	"github.com/phuslu/goid"
)

// IsInsideTest returns true inside a Go test
func IsInsideTest() bool {
	return len(os.Args) > 1 && strings.HasSuffix(os.Args[0], ".test") &&
		strings.HasPrefix(os.Args[1], "-test.")
}

// GoroutineID returns the current goroutine id.
func GoroutineID() int64 {
	return goid.Goid()
}
