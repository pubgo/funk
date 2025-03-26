package anyhow

import (
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

var errFnIsNil = errors.New("[fn] is nil")

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errFnIsNil)
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}

func tryResult[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errFnIsNil)
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return
}
