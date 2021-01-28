package xerror

import (
	"fmt"
)

func Assert(b bool, format string, a ...interface{}) {
	if !b {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fmt.Sprintf(format, a...) }))
}

func AssertFn(b bool, fn func() string) {
	if !b {
		return
	}

	panic(handle(ErrAssert, func(err *xerror) { err.Msg = fn() }))
}

func If(a bool, b interface{}, c interface{}) interface{} {
	if a {
		return b
	}

	return c
}
