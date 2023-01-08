package try

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
)

func WithErr(gErr *error, fn func() error) {
	assert.If(fn == nil, "[fn] is nil")

	var err errors.XError

	defer func() {
		if val := recover(); val != nil {
			err = errors.Parse(val)
		}

		if err != nil {
			err.AddTag("fn_caller", stack.CallerWithFunc(fn))
			*gErr = err
		}
	}()

	err = errors.Parse(fn())
}

func Try(fn func() error) (gErr error) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err = errors.Parse(val)
		if err == nil {
			return
		}

		err.AddTag("fn_caller", stack.CallerWithFunc(fn))
		gErr = err
	}()

	return fn()
}

func Result[T any](fn func() result.Result[T]) (g result.Result[T]) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		var err errors.XError
		if val := recover(); val != nil {
			err = errors.Parse(val)
		}

		if err == nil {
			return
		}

		err.AddTag("fn", stack.CallerWithFunc(fn))
		g = result.Err[T](err)
	}()

	return fn()
}
