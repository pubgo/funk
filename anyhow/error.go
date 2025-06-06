package anyhow

import (
	"context"

	"github.com/pubgo/funk/errors"
)

func newError(err error) Error {
	return Error{err: err}
}

type Error struct {
	_ [0]func() // disallow ==

	err error
}

func (r Error) WithErr(callbacks ...func(err error) error) Error {
	if r.IsOK() {
		return r
	}

	var err = r.getErr()
	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return Error{}
		}
	}
	return Error{err: errors.WrapCaller(err, 1)}
}

func (r Error) OnErr(callbacks ...func(err error)) {
	if r.IsOK() {
		return
	}

	var err = r.getErr()
	for _, fn := range callbacks {
		fn(err)
	}
}

func (r Error) CatchErr(setter *error, ctx ...context.Context) bool {
	return catchErr(r, nil, setter, ctx...)
}

func (r Error) Catch(setter *Error, ctx ...context.Context) bool {
	return catchErr(r, setter, nil, ctx...)
}

func (r Error) IsErr() bool { return r.getErr() != nil }

func (r Error) IsOK() bool { return r.getErr() == nil }

func (r Error) GetErr() error {
	if r.IsOK() {
		return nil
	}

	return errors.WrapCaller(r.getErr(), 1)
}

func (r Error) Must() {
	if r.IsOK() {
		return
	}

	errMust(errors.WrapCaller(r.getErr(), 1))
}

func (r Error) Expect(format string, args ...any) {
	if r.IsOK() {
		return
	}

	err := errors.WrapCaller(r.getErr(), 1)
	err = errors.Wrapf(err, format, args...)
	errMust(err)
}

func (r Error) getErr() error { return r.err }
