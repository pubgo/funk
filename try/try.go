package try

import (
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
)

func WithErr(gErr *error, fn func() error) {
	if fn == nil {
		*gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !errors.IsNil(err) {
			*gErr = errors.WrapStack(err)
		}

		*gErr = errors.WrapFn(*gErr, func(xrr errors.XError) {
			xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
		})
	}()

	*gErr = fn()
}

func Try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !errors.IsNil(err) {
			gErr = errors.WrapStack(err)
		}

		gErr = errors.WrapFn(gErr, func(xrr errors.XError) {
			xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
		})
	}()

	gErr = fn()
	return
}

func Result[T any](fn func() result.Result[T]) (g result.Result[T]) {
	if fn == nil {
		g = g.WithErr(errors.WrapStack(errors.New("[fn] is nil")))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !errors.IsNil(err) {
			g = g.WithErr(errors.WrapStack(err))
		}

		if g.IsErr() {
			g = g.WithErr(errors.WrapFn(g.Err(), func(xrr errors.XError) {
				xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
			}))
		}
	}()

	g = fn()
	return
}
