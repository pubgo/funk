package assert

import (
	"fmt"
)

var EnablePrintStack bool

func Assert(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func If(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func T(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func Err(b bool, err error) {
	if b {
		must(err)
	}
}

func Fn(b bool, fn func() error) {
	if b {
		must(fn())
	}
}

func Lazy(lazy func() bool, err error) {
	if lazy() {
		must(err)
	}
}
