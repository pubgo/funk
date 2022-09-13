package recovery

import (
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
)

func Result[T any](ret *result.Result[T]) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	*ret = result.Err[T](result.WithErr(err))
}

func ResultErr(gErr *result.Error, fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		*gErr = result.WithErr(fns[0](err1))
		return
	}
	*gErr = result.WithErr(err1)
}

func Err(gErr *error, fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		*gErr = fns[0](err1)
		return
	}
	*gErr = err1
}

func Raise(fns ...func(err xerr.XErr) xerr.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	err1 := xerr.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		panic(fns[0](err1))
	}
	panic(err1)
}

func Recovery(fn func(err xerr.XErr)) {
	assert.If(fn == nil, "[fn] should not be nil")

	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	fn(xerr.WrapXErr(err))
}

func Exit(handlers ...func()) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	if len(handlers) > 0 {
		handlers[0]()
	}
	xerr.WrapXErr(err).DebugPrint()
	os.Exit(1)
}

func DebugPrint() {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	xerr.WrapXErr(err).DebugPrint()
	debug.PrintStack()
}
