package anyhow

import (
	"runtime/debug"

	"github.com/pubgo/funk/anyhow/aherrcheck"
	"github.com/pubgo/funk/errors"
)

func newError(err error) Error {
	return Error{err: err}
}

type Error struct {
	err error
}

func (r Error) OnErr(callbacks ...func(err error)) Error {
	if r.IsErr() {
		var err = r.getErr()
		for _, fn := range callbacks {
			fn(err)
		}
	}

	return r
}

func (r Error) ErrTo(setter *Error, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	if !r.IsErr() {
		return
	}

	callbacks = append(callbacks, aherrcheck.GetErrChecks()...)
	var err = r.getErr()
	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	err = errors.WrapCaller(err, 1)
	*setter = newError(err)
}

func (r Error) IsErr() bool {
	return r.getErr() != nil
}

func (r Error) GetErr() error {
	if !r.IsErr() {
		return nil
	}

	return errors.WrapCaller(r.getErr(), 1)
}

func (r Error) getErr() error {
	return r.err
}

func (r Error) Unwrap(callback ...func(err error) error) {
	if !r.IsErr() {
		return
	}

	var err = errors.WrapCaller(r.getErr(), 1)
	for _, fn := range callback {
		err = fn(err)
		if err == nil {
			return
		}
	}

	debug.PrintStack()
	panic(err)
}

func (r Error) Expect(format string, args ...any) {
	if !r.IsErr() {
		return
	}

	debug.PrintStack()
	err := errors.WrapCaller(r.getErr(), 1)
	panic(errors.Wrapf(err, format, args...))
}
