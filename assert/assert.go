package assert

import (
	"fmt"

	"github.com/pubgo/funk/xerr"
)

func Assert(b bool, format string, a ...interface{}) {
	if b {
		panic(xerr.WrapXErr(fmt.Errorf(format, a...)))
	}
}

func If(b bool, format string, a ...interface{}) {
	if b {
		panic(xerr.WrapXErr(fmt.Errorf(format, a...)))
	}
}

func T(b bool, format string, a ...interface{}) {
	if b {
		panic(xerr.WrapXErr(fmt.Errorf(format, a...)))
	}
}

func Err(b bool, err error) {
	if b {
		panic(xerr.WrapXErr(err))
	}
}

func Fn(b bool, fn func() error) {
	if b {
		panic(xerr.WrapXErr(fn()))
	}
}

func Lazy(lazy func() bool, err error) {
	if lazy() {
		panic(xerr.WrapXErr(err))
	}
}