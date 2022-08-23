package result

import (
	"reflect"

	"github.com/pubgo/funk/xerr"
)

func WithErr(err error) Error {
	switch err.(type) {
	case nil:
		return Error{}
	case Error:
		return err.(Error)
	default:
		return Error{e: err}
	}
}

type Error struct {
	e error
}

func (e Error) IsNil() bool {
	return e.e == nil || reflect.ValueOf(e.e).IsNil()
}

func (e Error) Must() {
	if e.IsNil() {
		return
	}
	panic(xerr.Wrap(e.e))
}

func (e Error) Wrap(args ...interface{}) Error {
	if e.IsNil() {
		return e
	}
	return Error{e: xerr.Wrap(e.e, args...)}
}

func (e Error) Do(fn func(err Error)) {
	if e.IsNil() {
		return
	}

	fn(e)
}

func (e Error) WrapF(msg string, args ...interface{}) Error {
	if e.IsNil() {
		return e
	}
	return Error{e: xerr.WrapF(e.e, msg, args...)}
}

func (e Error) OrElse(wrap func(e Error) Error) Error {
	if e.IsNil() {
		return e
	}
	return wrap(e)
}

func (e Error) WithErr(err error) Error { e.e = err; return e }
func (e Error) Err() error              { return e.e }
func (e Error) Unwrap() error           { return e.e }
func (e Error) Error() string {
	if e.IsNil() {
		return ""
	}
	return e.e.Error()
}
func (e Error) Expect(msg string, args ...interface{}) {
	if e.IsNil() {
		return
	}
	panic(xerr.WrapF(e.e, msg, args...))
}
