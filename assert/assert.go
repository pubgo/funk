package assert

import (
	"fmt"
	"github.com/pubgo/funk/errors"
)

func Assert(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapStack(fmt.Errorf(format, a...)))
	}
}

func If(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapStack(fmt.Errorf(format, a...)))
	}
}

func T(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapStack(fmt.Errorf(format, a...)))
	}
}

func Err(b bool, err error) {
	if b {
		panic(errors.WrapStack(err))
	}
}

func Fn(b bool, fn func() error) {
	if b {
		panic(errors.WrapStack(fn()))
	}
}

func Lazy(lazy func() bool, err error) {
	if lazy() {
		panic(errors.WrapStack(err))
	}
}
