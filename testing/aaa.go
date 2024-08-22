package testing

type TestingTB interface {
	Name() string
	Cleanup(f func())
	Logf(fmt string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(message string, args ...interface{})
}
