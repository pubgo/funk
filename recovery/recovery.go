package recovery

import (
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/result"
)

func Result[T any](ret *result.Result[T]) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	*ret = result.Err[T](errors.WrapStack(err))
}

func Err(gErr *error, fns ...func(err error) error) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	for i := range fns {
		err = fns[i](err)
	}

	*gErr = errors.WrapStack(err)
}

func Raise(fns ...func(err error) error) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	if len(fns) > 0 && fns[0] != nil {
		panic(errors.WrapCaller(fns[0](err)))
	}

	panic(errors.WrapStack(err))
}

func Recovery(fn func(err error)) {
	assert.If(fn == nil, "[fn] should not be nil")

	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	fn(errors.WrapStack(err))
}

func Exit(handlers ...func(err error) error) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	for i := range handlers {
		err = handlers[i](err)
	}

	errors.Debug(errors.WrapStack(err))
	debug.PrintStack()
	os.Exit(1)
}

func DebugPrint() {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapStack(err))
	debug.PrintStack()
}
