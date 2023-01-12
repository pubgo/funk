package try

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
)

func WithErr(gErr *error, fn func() error) {
	assert.If(fn == nil, "[fn] is nil")

	var err error
	defer func() {
		if err1 := errors.Parse(recover()); !errors.IsNil(err1) {
			err = errors.WrapStack(err1)
		}

		*gErr = errors.WrapFn(err, func(xrr errors.XError) {
			xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
		})
	}()

	err = fn()
}

func Try(fn func() error) (gErr error) {
	assert.If(fn == nil, "[fn] is nil")

	var err error
	defer func() {
		if err1 := errors.Parse(recover()); !errors.IsNil(err1) {
			err = errors.WrapStack(err1)
		}

		gErr = errors.WrapFn(err, func(xrr errors.XError) {
			xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
		})
	}()
	err = fn()
	return
}

func Result[T any](fn func() result.Result[T]) (g result.Result[T]) {
	assert.If(fn == nil, "[fn] is nil")

	var err result.Result[T]
	defer func() {
		if err1 := errors.Parse(recover()); !errors.IsNil(err1) {
			err = err.WithErr(errors.WrapStack(err1))
		}

		if err.IsErr() {
			err = err.WithErr(errors.WrapFn(err.Err(), func(xrr errors.XError) {
				xrr.AddTag("fn_stack", stack.CallerWithFunc(fn))
			}))
		}

		g = err
	}()

	err = fn()
	return
}
