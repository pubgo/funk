package assert

import (
	"github.com/pubgo/funk/errors"
)

func Assert(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapCaller(errors.New(format, a...), 1))
	}
}

func If(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapCaller(errors.New(format, a...), 1))
	}
}

func T(b bool, format string, a ...interface{}) {
	if b {
		panic(errors.WrapCaller(errors.New(format, a...), 1))
	}
}

func Err(b bool, err error) {
	if b {
		panic(errors.WrapCaller(err, 1))
	}
}

func Fn(b bool, fn func() error) {
	if b {
		panic(errors.WrapCaller(fn(), 1))
	}
}

func Lazy(lazy func() bool, err error) {
	if lazy() {
		panic(errors.WrapCaller(err, 1))
	}
}
