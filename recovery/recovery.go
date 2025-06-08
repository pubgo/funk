package recovery

import (
	"os"
	"testing"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
)

func Err(gErr *error, callbacks ...func(err error) error) {
	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	for i := range callbacks {
		err = callbacks[i](err)
		if funk.IsNil(err) {
			return
		}
	}

	*gErr = errors.WrapStack(err)
}

func Raise(callbacks ...func(err error) error) {
	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	for i := range callbacks {
		err = callbacks[i](err)
		if funk.IsNil(err) {
			return
		}
	}

	panic(errors.WrapStack(err))
}

func Recovery(fn func(err error)) {
	assert.If(fn == nil, "[fn] should not be nil")

	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	fn(errors.WrapStack(err))
}

func Exit(handlers ...func(err error) error) {
	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	for i := range handlers {
		err = handlers[i](err)
		if funk.IsNil(err) {
			return
		}
	}

	errors.Debug(errors.WrapStack(err))
	os.Exit(1)
}

func DebugPrint() {
	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapStack(err))
}

func Testing(t *testing.T) {
	err := errors.Parse(recover())
	if funk.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapStack(err))
	t.Fatal(err)
}
