package xtry

import (
	"fmt"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
)

func TryWith(gErr *result.Error, fn func() result.Error) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		xerr.ParseErr(&err, val)
		if err == nil {
			return
		}

		*gErr = result.WithErr(xerr.WrapF(err, "fn=%s", utils.CallerWithFunc(fn)))
	}()

	*gErr = fn().OrElse(func(e result.Error) result.Error {
		return e.WrapF("fn=%s", utils.CallerWithFunc(fn))
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
		xerr.ParseErr(&err, val)
		if err == nil {
			return
		}

		gErr = result.WithErr(err).OrElse(func(e result.Error) result.Error {
			return e.WrapF("fn=%s", utils.CallerWithFunc(fn))
		})
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
		xerr.ParseErr(&err, val)
		if err == nil {
			return
		}

		gErr = result.WithErr(xerr.Wrap(err, fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))))
	}()

	return fn().OrElse(func(e result.Error) result.Error {
		return e.WrapF("fn=%s", utils.CallerWithFunc(fn))
	})
}

func TryCatch(fn func() result.Error, catch func(err result.Error)) {
	assert.If(fn == nil, "[fn] is nil")
	assert.If(catch == nil, "[catch] is nil")

	var gErr result.Error
	defer func() {
		if !gErr.IsNil() {
			catch(gErr.WrapF("fn=%s", utils.CallerWithFunc(fn)))
			return
		}

		val := recover()
		if val == nil {
			return
		}

		var err error
		xerr.ParseErr(&err, val)
		if err == nil {
			return
		}

		catch(result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn)))
	}()

	gErr = fn()
}

func TryCatch1[T any](fn func() result.Result[T]) (g result.Result[T]) {
	assert.If(fn == nil, "[fn] is nil")

	defer func() {
		val := recover()
		if val == nil {
			return
		}

		var err error
		xerr.ParseErr(&err, val)
		if err == nil {
			return
		}

		g = result.Err[T](result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn)))
	}()

	return fn()
}
