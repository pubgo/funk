package result

import (
	"fmt"

	"github.com/pubgo/funk/errors"
)

func WithErr(err error) Error {
	return Error{e: err}
}

func NilErr() Error {
	return Error{}
}

type Error struct {
	e error
}

func (e Error) String() string {
	if e.IsNil() {
		return ""
	}

	return fmt.Sprintf("err=%q detail=%#v", e.e.Error(), e.e)
}

func (e Error) IsNil() bool { return errors.IsNil(e.e) }

func (e Error) Wrap(args ...interface{}) Error {
	if e.IsNil() {
		return e
	}

	return Error{e: errors.Wrap(e.e, fmt.Sprint(args...))}
}

func (e Error) Do(fn func(err Error)) {
	if e.IsNil() {
		return
	}

	fn(e)
}

func (e Error) Wrapf(msg string, args ...interface{}) Error {
	if e.IsNil() {
		return e
	}

	return Error{e: errors.Wrapf(e.e, msg, args...)}
}

func (e Error) OrElse(wrap func(e Error) Error) Error {
	if e.IsNil() {
		return e
	}

	return wrap(e)
}

func (e Error) WithErr(err error) Error { e.e = err; return e }
func (e Error) WithTag(k string, v interface{}) Error {
	if e.IsNil() {
		return e
	}

	e.e = errors.WrapFn(e.e, func(xrr errors.XError) {
		xrr.AddTag(k, v)
	})
	return e
}
func (e Error) Err() error    { return e.e }
func (e Error) Unwrap() error { return e.e }
func (e Error) Expect(msg string, args ...interface{}) {
	if e.IsNil() {
		return
	}

	panic(errors.Wrapf(e.e, msg, args...))
}
