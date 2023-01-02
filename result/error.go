package result

import (
	"fmt"

	"github.com/pubgo/funk/errors"
)

func ParseErr(err *error, val interface{}) {
	switch _val := val.(type) {
	case nil:
		return
	case Error:
		*err = _val.Err()
	case interface{ Unwrap() error }:
		if _val.Unwrap() != nil {
			*err = errors.WrapXErr(_val.Unwrap()).WrapF("%#v", _val)
		} else {
			*err = fmt.Errorf("%#v", _val)
		}
	case error:
		*err = _val
	case string:
		*err = errors.New(_val)
	case []byte:
		*err = errors.New(string(_val))
	default:
		*err = fmt.Errorf("%#v", _val)
	}
	*err = errors.WrapXErr(*err)
}

func WithErr(err error) Error {
	switch err.(type) {
	case nil:
		return Error{}
	default:
		return Error{e: errors.WrapXErr(err)}
	}
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

func (e Error) IsNil() bool { return errors.isNil(e.e) }

func (e Error) IsErr() bool { return !e.IsNil() }

func (e Error) Must(check ...func(err Error) Error) {
	if e.IsNil() {
		return
	}

	if len(check) > 0 && check[0] != nil {
		panic(errors.Wrap(check[0](e).e))
	} else {
		panic(errors.Wrap(e.e))
	}
}

func (e Error) Wrap(args ...interface{}) Error {
	if e.IsNil() {
		return e
	}

	return Error{e: errors.Wrap(e.e, args...)}
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

	return Error{e: errors.WrapF(e.e, msg, args...)}
}

func (e Error) OrElse(wrap func(e Error) Error) Error {
	if e.IsNil() {
		return e
	}

	return wrap(e)
}

func (e Error) WithErr(err error) Error { e.e = err; return e }
func (e Error) WithMeta(k string, v interface{}) Error {
	if e.IsNil() {
		return e
	}

	e.e = errors.WrapXErr(e.e).WithMeta(k, v)
	return e
}
func (e Error) Err() error    { return e.e }
func (e Error) Unwrap() error { return e.e }
func (e Error) Expect(msg string, args ...interface{}) {
	if e.IsNil() {
		return
	}

	panic(errors.WrapF(e.e, msg, args...))
}
