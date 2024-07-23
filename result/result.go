package result

import (
	"encoding/json"
	"fmt"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

var _ error = (*Error)(nil)

type Error struct {
	Msg   string
	Stack string
}

func (e Error) String() string {
	return fmt.Sprintf("%s:\n%s", e.Msg, e.Stack)
}

func (e Error) Error() string {
	return e.String()
}

type R[T any] interface {
	Unwrap() T
	IsErr() bool
	Err() error
	Expect(format string, args ...any) T
}

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v}
}

func Err[T any](err error) Result[T] {
	return Result[T]{e: errors.WrapCaller(err, 1)}
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, e: errors.WrapCaller(err, 1)}
}

func Of[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, e: errors.WrapCaller(err, 1)}
}

type Result[T any] struct {
	v *T
	e error
}

func (r Result[T]) WithErr(err error) Result[T] {
	return Result[T]{e: errors.WrapCaller(err, 1)}
}

func (r Result[T]) WithVal(v T) Result[T] {
	return OK(v)
}

func (r Result[T]) ValueTo(v *T) error {
	if r.IsErr() {
		return errors.WrapCaller(r.e, 1)
	}

	*v = generic.FromPtr(r.v)
	return nil
}

func (r Result[T]) OnValue(fn func(t T) error) error {
	if r.IsErr() {
		return r.e
	}

	return errors.WrapCaller(fn(generic.FromPtr(r.v)), 1)
}

func (r Result[T]) Err(check ...func(err error) error) error {
	if !r.IsErr() {
		return nil
	}

	if len(check) > 0 && check[0] != nil {
		return errors.WrapCaller(check[0](r.e), 1)
	}

	return errors.WrapCaller(r.e, 1)
}

func (r Result[T]) IsErr() bool {
	return r.e != nil || !generic.IsNil(r.e)
}

func (r Result[T]) OrElse(v T) T {
	if r.IsErr() {
		return v
	}
	return generic.FromPtr(r.v)
}

func (r Result[T]) Unwrap(check ...func(err error) error) T {
	if !r.IsErr() {
		return generic.FromPtr(r.v)
	}

	if len(check) > 0 && check[0] != nil {
		panic(check[0](r.e))
	} else {
		panic(r.e)
	}
}

func (r Result[T]) Expect(format string, args ...any) T {
	if !r.IsErr() {
		return generic.FromPtr(r.v)
	}

	panic(&Error{
		Msg:   fmt.Sprintf(format, args...),
		Stack: stack.Caller(1).String(),
	})
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", generic.FromPtr(r.v))
	}

	return fmt.Sprint(errors.WrapCaller(r.e, 1))
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, errors.WrapCaller(r.e, 1)
	}

	return json.Marshal(generic.FromPtr(r.v))
}

func (r Result[T]) UnmarshalJSON([]byte) error {
	panic("unimplemented")
}
