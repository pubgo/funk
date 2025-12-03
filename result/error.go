package result

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errutil"
	"github.com/pubgo/funk/v2/log/logfields"
)

var (
	_ Catchable = new(Error)
	_ Checkable = new(Error)
	_ ErrSetter = new(Error)
)

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

func (e Error) LogCtx(ctx context.Context, events ...func(e *zerolog.Event)) Error {
	logErr(ctx, 0, e.err, events...)
	return e
}

func (e Error) Log(events ...func(e *zerolog.Event)) Error {
	logErr(context.Background(), 0, e.err, events...)
	return e
}

func (e Error) WrapErr(err *errors.Err, tags ...errors.Tag) Error {
	return Error{err: errors.WrapTag(errors.WrapCaller(err, 1), tags...)}
}

func (e Error) WithFn(fn func() error) Error {
	if fn == nil {
		return Error{err: errors.WrapCaller(errFnIsNil, 1)}
	}
	err := fn()
	if err == nil {
		return Error{}
	}
	return Error{err: errors.WrapCaller(err, 1)}
}

func (e Error) WithErr(err error) Error {
	return Error{err: errors.WrapCaller(err, 1)}
}

func (e Error) WithErrorf(format string, args ...any) Error {
	return Error{err: errors.WrapCaller(fmt.Errorf(format, args...), 1)}
}

func (e Error) Inspect(fn func(error)) Error {
	if e.IsErr() {
		err := e.getErr()
		fn(err)
	}

	return e
}

func (e Error) InspectErr(fn func(error)) Error { return e.Inspect(fn) }

func (e Error) CatchErr(setter *error, ctx ...context.Context) bool {
	return catchErr(e, nil, setter, ctx...)
}

func (e Error) Catch(setter ErrSetter, ctx ...context.Context) bool {
	return catchErr(e, setter, nil, ctx...)
}

func (e Error) IsErr() bool { return e.getErr() != nil }

func (e Error) IsOK() bool { return e.getErr() == nil }

func (e Error) GetErr() error {
	if e.IsOK() {
		return nil
	}

	return e.getErr()
}

func (e Error) Must(events ...func(e *zerolog.Event)) {
	if e.IsOK() {
		return
	}

	errNilOrPanic(errors.WrapCaller(e.getErr(), 1), events...)
}

func (e Error) Expect(format string, args ...any) {
	if e.IsOK() {
		return
	}

	err := errors.WrapCaller(e.getErr(), 1)
	errNilOrPanic(err, func(e *zerolog.Event) {
		e.Str(logfields.Msg, fmt.Sprintf(format, args...))
	})
}

func (e Error) String() string {
	if e.IsOK() {
		return "OK"
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

func (e Error) setErrorInner() {
}
