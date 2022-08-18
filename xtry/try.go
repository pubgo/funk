package funk

import (
	"fmt"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/xerr"
)

func TryWith(gErr *error, fn func() error) {
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

		*gErr = xerr.Wrap(err, fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn)))
	}()

	*gErr = xerr.WrapXErr(fn(), func(err *xerr.XError) {
		err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
	})
}

func Try(fn func() error) (gErr error) {
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

		gErr = xerr.Wrap(err, fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn)))
	}()

	return xerr.WrapXErr(fn(), func(err *xerr.XError) {
		err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
	})
}

func TryCatch(fn func() error, catch func(err xerr.XErr)) {
	assert.If(fn == nil, "[fn] is nil")
	assert.If(catch == nil, "[catch] is nil")

	var gErr error
	defer func() {
		if gErr != nil {
			catch(xerr.WrapXErr(gErr, func(err *xerr.XError) {
				err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
			}))
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

		catch(xerr.WrapXErr(err, func(err *xerr.XError) {
			err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
		}))
	}()

	gErr = fn()
}

func TryCatch1[T any](fn func() (T, error), catch func(err xerr.XErr)) (v T) {
	assert.If(fn == nil, "[fn] is nil")
	assert.If(catch == nil, "[catch] is nil")

	var gErr error
	defer func() {
		if gErr != nil {
			catch(xerr.WrapXErr(gErr, func(err *xerr.XError) {
				err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
			}))
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

		catch(xerr.WrapXErr(err, func(err *xerr.XError) {
			err.Detail = fmt.Sprintf("fn=%s", utils.CallerWithFunc(fn))
		}))
	}()

	v, gErr = fn()
	return v
}
