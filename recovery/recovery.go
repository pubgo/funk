package recovery

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
)

func Err(gErr *error, callbacks ...func(err error) error) {
	err := errparser.Parse(recover())
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
	err := errparser.Parse(recover())
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

	err := errparser.Parse(recover())
	if err == nil {
		return
	}

	debug.PrintStack()
	fn(err)
}

func Exit(handlers ...func(err error) error) {
	err := errparser.Parse(recover())
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
	errors.DebugPrint(err)
	os.Exit(1)
}

func DebugPrint() {
	err := errparser.Parse(recover())
	if err == nil {
		return
	}

	debug.PrintStack()
	errors.DebugPrint(err)
}

func Testing(t *testing.T) {
	err := errparser.Parse(recover())
	if err == nil {
		return
	}

	errors.DebugPrint(err)
	t.Fatal(err)
}
