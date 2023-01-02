package recovery

import (
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
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

	debug.PrintStack()
	*ret = result.Err[T](err)
}

func ResultErr(gErr *result.Error, fns ...func(err errors.XErr) errors.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	debug.PrintStack()
	err1 := errors.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		*gErr = result.WithErr(fns[0](err1))
		return
	}
	*gErr = result.WithErr(err1)
}

func Err(gErr *error, fns ...func(err errors.XErr) errors.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	debug.PrintStack()
	err1 := errors.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		*gErr = fns[0](err1)
		return
	}
	*gErr = err1
}

func Raise(fns ...func(err errors.XErr) errors.XErr) {
	val := recover()
	if val == nil {
		return
	}

	var err error
	result.ParseErr(&err, val)
	if err == nil {
		return
	}

	debug.PrintStack()
	err1 := errors.WrapXErr(err)
	if len(fns) > 0 && fns[0] != nil {
		panic(fns[0](err1))
	}
	panic(err1)
}

func Recovery(fn func(err errors.XErr)) {
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

	debug.PrintStack()
	fn(errors.WrapXErr(err))
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
	debug.PrintStack()
	errors.WrapXErr(err).DebugPrint()
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

	debug.PrintStack()
	errors.WrapXErr(err).DebugPrint()
}
