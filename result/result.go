package result

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/samber/lo"
)

type ErrSetter interface {
	setErr(err error)
}

type R[T any] interface {
	Unwrap() T
	IsErr() bool
	Err() error
	Expect(format string, args ...any) T
}

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v, e: newError(nil)}
}

func Err[T any](err error) Result[T] {
	err = errors.WrapCaller(err, 1)
	return Result[T]{e: newError(err)}
}

func newError(err error) Error {
	return Error{err: err}
}

func ErrOf(err error) Error {
	err = errors.WrapCaller(err, 1)
	return newError(err)
}

func Wrap[T any](v T, err error) Result[T] {
	if err == nil {
		return Result[T]{v: &v, e: newError(nil)}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{e: newError(err)}
}

func Of[T any](v T, err error) Result[T] {
	if err == nil {
		return Result[T]{v: &v, e: newError(nil)}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{e: newError(err)}
}

type Error struct {
	err error
}

func (r *Error) setErr(err error) {
	if err == nil {
		return
	}

	r.err = err
}

func (r Error) OnErr(fn func(err error)) {
	if r.err == nil {
		return
	}

	fn(r.err)
}

func (r Error) ErrTo(setter ErrSetter, callback ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	if r.err == nil {
		return
	}

	callback = append(callback, errChecks...)
	var err = errors.WrapCaller(r.err, 1)
	for _, fn := range callback {
		err = fn(err)
		if err == nil {
			return
		}
	}
	setter.setErr(err)
}

func (r Error) IsErr() bool {
	return r.err != nil
}

func (r Error) Err() error {
	if !r.IsErr() {
		return nil
	}

	return errors.WrapCaller(r.err, 1)
}

type Result[T any] struct {
	e Error
	v *T
}

func (r Result[T]) WithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{e: newError(err)}
}

func (r Result[T]) WithVal(v T) Result[T] {
	return OK(v)
}

func (r Result[T]) ValueTo(v *T) error {
	if r.e.IsErr() {
		return errors.WrapCaller(r.e.Err(), 1)
	}

	*v = lo.FromPtr(r.v)
	return nil
}

func (r Result[T]) OnValue(fn func(t T) error) error {
	if r.e.IsErr() {
		return errors.WrapCaller(r.e.Err(), 1)
	}

	return errors.WrapCaller(fn(lo.FromPtr(r.v)), 1)
}

func (r Result[T]) OrElse(v T, callback ...func(err error)) T {
	var isErr = r.e.IsErr()
	if isErr {
		for _, fn := range callback {
			fn(r.e.Err())
		}
	}

	return lo.If(isErr, v).ElseF(func() T {
		return lo.FromPtr(r.v)
	})
}

func (r Result[T]) Unwrap(check ...func(err error) error) T {
	if !r.e.IsErr() {
		return lo.FromPtr(r.v)
	}

	var err = r.e.err
	for _, fn := range check {
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
	panic(errors.Wrapf(r.e.Err(), format, args...))
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", lo.FromPtr(r.v))
	}

	return fmt.Sprint(errors.WrapCaller(r.Err(), 1))
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, errors.WrapCaller(r.Err(), 1)
	}

	return json.Marshal(lo.FromPtr(r.v))
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	if r.IsErr() {
		return errors.WrapCaller(r.Err(), 1)
	}

	var v T
	var err = json.Unmarshal(data, &v)
	if err != nil {
		return errors.WrapCaller(err, 1)
	}
	r.v = &v
	return nil
}

func (r Result[T]) Do(fn func(v T)) {
	if r.IsErr() {
		return
	}

	fn(lo.FromPtr(r.v))
}

func (r Result[T]) ErrTo(setter ErrSetter, callback ...func(err error) error) T {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	var ret = lo.FromPtr(r.v)
	if r.IsErr() {
		callback = append(callback, errChecks...)
		var err = errors.WrapCaller(r.e.Err(), 1)
		for _, fn := range callback {
			err = fn(err)
			if err == nil {
				return ret
			}
		}
		setter.setErr(err)
	}
	return ret
}

func (r Result[T]) IsErr() bool {
	return r.e.IsErr()
}

func (r Result[T]) Err(callback ...func(err error) error) error {
	if !r.IsErr() {
		return nil
	}

	var err = r.e.Err()
	for _, cc := range callback {
		err = cc(err)
		if err == nil {
			return nil
		}
	}

	return errors.WrapCaller(err, 1)
}

func (r *Result[T]) setErr(err error) {
	if err == nil {
		return
	}

	r.e = newError(err)
}
