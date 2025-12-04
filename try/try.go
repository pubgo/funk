package try

import (
	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
	"github.com/pubgo/funk/v2/stack"
)

func WithErr(gErr *error, fn func() error) {
	if fn == nil {
		*gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return
	}

	defer func() {
		if err := errparser.Parse(recover()); !funk.IsNil(err) {
			*gErr = errors.WrapStack(err)
		}

		*gErr = errors.WrapKV(*gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	*gErr = fn()
}

func Try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errors.New("[fn] is nil"))
		return gErr
	}

	defer func() {
		if err := errparser.Parse(recover()); !funk.IsNil(err) {
			gErr = errors.WrapStack(err)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return gErr
}
