package testing

import "testing"

var _ Interface = (*testing.T)(nil)
var _ Interface = (*testing.B)(nil)

type Interface interface {
	Name() string
	Cleanup(f func())
	Logf(fmt string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(message string, args ...interface{})
}
