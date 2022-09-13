package xtry

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/result"
)

func TryWith(gErr *result.Error, fn func() result.Error) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		result.ParseErr(&err, val)
		if err == nil {
			return
		}

		*gErr = result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn))
	}()

	fn().Do(func(err result.Error) {
		*gErr = err.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})
}

func Try(fn func()) (gErr result.Error) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		result.ParseErr(&err, val)
		if err == nil {
			return
		}

		gErr = result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn))
	}()

	fn()
	return
}

func TryErr(fn func() result.Error) (gErr result.Error) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		result.ParseErr(&err, val)
		if err == nil {
			return
		}

		gErr = result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn))
	}()

	return fn().OrElse(func(e result.Error) result.Error {
		return e.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})
}

func TryVal[T any](fn func() result.Result[T]) (g result.Result[T]) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		result.ParseErr(&err, val)
		if err == nil {
			return
		}

		g = result.Err[T](result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn)))
	}()

	return fn()
}
