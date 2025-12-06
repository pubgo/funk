package testing

import "testing"

var (
	_ Interface = (*testing.T)(nil)
	_ Interface = (*testing.B)(nil)
)

type Interface interface {
	Name() string
	Cleanup(f func())
	Logf(fmt string, args ...any)
	Fatalf(format string, args ...any)
	Errorf(message string, args ...any)
}
