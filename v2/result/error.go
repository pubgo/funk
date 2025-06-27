package result

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/errors/errutil"
	"github.com/pubgo/funk/log"
)

var _ Catchable = new(Error)
var _ Checkable = new(Error)
var _ ErrSetter = new(Error)

func newError(err error) Error {
	return Error{err: err}
}

type Error struct {
	_ [0]func() // disallow ==

	err error
}

func (e Error) Map(fn func(error) error) Error {
	if e.IsOK() {
		return e
	}

	err := e.getErr()
	err = errors.WrapCaller(fn(err), 1)
	return Error{err: err}
}

func (e Error) LogErr(contexts ...context.Context) Error {
	if e.IsErr() {
		log.Err(e.err, contexts...).
			CallerSkipFrame(1).
			Msg(e.err.Error())
	}

	return e
}

func (e Error) WrapErr(err *errors.Err, tags ...errors.Tag) Error {
	return Error{err: errors.WrapTag(errors.WrapCaller(err, 1), tags...)}
}

func (e Error) Inspect(fn func(error)) Error {
	if e.IsErr() {
		err := e.getErr()
		fn(err)
	}

	return e
}

func (e Error) Unwrap() error { return e.err }

func (e Error) Catch(setter *error, ctx ...context.Context) bool {
	return catchErr(e, nil, setter, ctx...)
}

func (e Error) CatchErr(setter ErrSetter, ctx ...context.Context) bool {
	return catchErr(e, setter, nil, ctx...)
}

func (e Error) IsErr() bool { return e.getErr() != nil }

func (e Error) IsOK() bool { return e.getErr() == nil }

func (e Error) GetErr() error {
	if e.IsOK() {
		return nil
	}

	return errors.WrapCaller(e.getErr(), 1)
}

func (e Error) Must() {
	if e.IsOK() {
		return
	}

	errMust(errors.WrapCaller(e.getErr(), 1))
}

func (e Error) Expect(format string, args ...any) {
	if e.IsOK() {
		return
	}

	err := errors.WrapCaller(e.getErr(), 1)
	err = errors.Wrapf(err, format, args...)
	errMust(err)
}

func (e Error) String() string {
	if e.IsOK() {
		return "Ok"
	}

	return fmt.Sprintf("Error(%v)", e.err)
}

func (e Error) MarshalJSON() ([]byte, error) {
	if e.IsErr() {
		return nil, errors.WrapCaller(e.err, 1)
	}

	return errutil.Json(e.err), nil
}

func (e Error) getErr() error { return e.err }

func (e *Error) setError(err error) {
	if err == nil {
		return
	}
	e.err = err
}
