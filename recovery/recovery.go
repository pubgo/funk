package recovery

import (
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
)

func Result[T any](ret *result.Result[T]) {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	*ret = result.Err[T](err)
}

func ResultErr(gErr *error, fns ...func(err errors.XError)) {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	if len(fns) > 0 && fns[0] != nil {
		fns[0](err)
	}

	*gErr = err
}

func Err(gErr *error, fns ...func(err errors.XError)) {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	if len(fns) > 0 && fns[0] != nil {
		fns[0](err)
		return
	}

	*gErr = err
}

func Raise(fns ...func(err errors.XError)) {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	debug.PrintStack()
	if len(fns) > 0 && fns[0] != nil {
		fns[0](err)
	}

	panic(errors.WrapCaller(err, 1))
}

func Recovery(fn func(err errors.XError)) {
	assert.If(fn == nil, "[fn] should not be nil")

	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	fn(err)
}

func Exit(handlers ...func()) {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	if len(handlers) > 0 {
		handlers[0]()
	}

	errors.Debug(err)
	debug.PrintStack()
	os.Exit(1)
}

func DebugPrint() {
	err := errors.Parse(recover())
	if errors.IsNil(err) {
		return
	}

	errors.Debug(err)
	debug.PrintStack()
}
