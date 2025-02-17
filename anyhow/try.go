package anyhow

import (
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}

func tryResult[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errors.New("[fn] is nil"))
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return
}
