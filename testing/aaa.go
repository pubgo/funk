// testing/testing.go

package testing

import "testing"

var (
	_ Interface = (*testing.T)(nil)
	_ Interface = (*testing.B)(nil)
)

// Interface is a testing interface.
type Interface interface {
	Name() string
	Cleanup(f func())
	Logf(fmt string, args ...any)
	Fatalf(format string, args ...any)
	Errorf(message string, args ...any)
}
