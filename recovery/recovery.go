package recovery

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/pubgo/funk/v2/errors/errinter"
	"github.com/samber/lo"
)

func Err(gErr *error, callbacks ...func(err error) error) {
	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	for i := range callbacks {
		err = callbacks[i](err)
		if err == nil {
			return
		}
	}

	debug.PrintStack()
	*gErr = err
}

func Raise(callbacks ...func(err error) error) {
	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	for i := range callbacks {
		err = callbacks[i](err)
		if err == nil {
			return
		}
	}

	debug.PrintStack()
	panic(err)
}

func Recovery(fn func(err error)) {
	lo.Assert(fn != nil, "[fn] should not be nil")

	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	debug.PrintStack()
	fn(err)
}

func Exit(handlers ...func(err error) error) {
	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	for i := range handlers {
		err = handlers[i](err)
		if err == nil {
			return
		}
	}

	debug.PrintStack()
	errinter.Debug(err)
	os.Exit(1)
}

func DebugPrint() {
	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	debug.PrintStack()
	errinter.Debug(err)
}

func Testing(t *testing.T) {
	err := errinter.ParseError(recover())
	if err == nil {
		return
	}

	errinter.Debug(err)
	t.Fatal(err)
}
