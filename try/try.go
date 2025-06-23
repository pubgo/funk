package try

import (
	"github.com/pubgo/funk"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
)

// Result
// Deprecated: use result.Try instead
func Result[T any](fn func() result.Result[T]) (g result.Result[T]) {
	if fn == nil {
		g = g.WithErr(errors.WrapStack(errors.New("[fn] is nil")))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			g = g.WithErr(errors.WrapStack(err))
		}

		if g.IsErr() {
			g = g.WithErr(errors.WrapKV(g.Err(), "fn_stack", stack.CallerWithFunc(fn)))
		}
	}()

	g = fn()
	return
}

func WithErr(gErr *error, fn func() error) {
	if fn == nil {
		*gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !funk.IsNil(err) {
			*gErr = errors.WrapStack(err)
		}

		*gErr = errors.WrapKV(*gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	*gErr = fn()
}

func Try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !funk.IsNil(err) {
			gErr = errors.WrapStack(err)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}
