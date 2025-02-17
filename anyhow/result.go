package anyhow

import (
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/samber/lo"
)

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		return Result[T]{v: nil}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
}

func newError(err error) Error {
	return Error{err: err}
}

func ErrOf(err error) Error {
	if err == nil {
		return Error{}
	}

	err = errors.WrapCaller(err, 1)
	return newError(err)
}

func ErrOfFn(fn func() error) Error {
	var err = try(fn)
	err = errors.WrapCaller(err, 1)
	return newError(err)
}

func Wrap[T any](v T, err error) Result[T] {
	if err == nil {
		return Result[T]{v: &v}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
}

func WrapFn[T any](fn func() (T, error)) Result[T] {
	v, err := tryResult(fn)
	if err == nil {
		return Result[T]{v: &v}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
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

	callbacks = append(callbacks, errChecks...)
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

	return errors.WrapCaller(r.err, 1)
}

func (r Error) getErr() error {
	return r.err
}

type Result[T any] struct {
	v   *T
	Err Error
}

func (r Result[T]) GetValue() T {
	return r.getValue()
}

func (r Result[T]) OnValue(fn func(v T)) Result[T] {
	if !r.IsErr() {
		fn(lo.FromPtr(r.v))
	}
	return r
}

func (r Result[T]) WithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
}

func (r Result[T]) WithVal(v T) Result[T] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Result[T]{Err: newError(err)}
	}

	return OK(v)
}

func (r Result[T]) ValueTo(v *T) error {
	if r.IsErr() {
		return errors.WrapCaller(r.getErr(), 1)
	}

	if v == nil {
		debug.PrintStack()
		panic("v params is nil")
	}

	*v = lo.FromPtr(r.v)
	return nil
}

func (r Result[T]) Unwrap(callback ...func(err error) error) T {
	if !r.Err.IsErr() {
		return lo.FromPtr(r.v)
	}

	var err = errors.WrapCaller(r.getErr(), 1)
	for _, fn := range callback {
		err = fn(err)
		if err == nil {
			return lo.FromPtr(r.v)
		}
	}

	debug.PrintStack()
	panic(err)
}

func (r Result[T]) Expect(format string, args ...any) T {
	if !r.IsErr() {
		return lo.FromPtr(r.v)
	}

	debug.PrintStack()
	err := errors.WrapCaller(r.getErr(), 1)
	panic(errors.Wrapf(err, format, args...))
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", lo.FromPtr(r.v))
	}

	return fmt.Sprint(errors.WrapCaller(r.getErr(), 1))
}

func (r Result[T]) ErrTo(setter *Error, callbacks ...func(err error) error) T {
	if setter == nil {
		debug.PrintStack()
		panic("ErrTo: setter is nil")
	}

	var ret = lo.FromPtr(r.v)
	if !r.IsErr() {
		return ret
	}

	callbacks = append(callbacks, errChecks...)
	var err = r.getErr()
	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return ret
		}
	}
	err = errors.WrapCaller(err, 1)
	*setter = newError(err)
	return ret
}

func (r Result[T]) IsErr() bool {
	return r.getErr() != nil
}

func (r Result[T]) GetErr() error {
	var err = r.getErr()
	return errors.WrapCaller(err, 1)
}

func (r Result[T]) OnErr(callbacks ...func(err error)) Result[T] {
	if r.IsErr() {
		var err = r.getErr()
		for _, fn := range callbacks {
			fn(err)
		}
	}

	return r
}

func (r Result[T]) getValue() T {
	return lo.FromPtr(r.v)
}

func (r Result[T]) getErr() error {
	return r.Err.getErr()
}
