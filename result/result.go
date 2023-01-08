package result

import (
	"encoding/json"
	"fmt"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
)

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v}
}

func Err[T any](err error) Result[T] {
	return Result[T]{e: err}
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, e: err}
}

type Result[T any] struct {
	v *T
	e error
}

func (r Result[T]) WithErr(err error) Result[T] {
	r.e = err
	return r
}

func (r Result[T]) WithVal(v T) Result[T] {
	r.v = &v
	return r
}

func (r Result[T]) Err(check ...func(err error) error) error {
	if !r.IsErr() {
		return nil
	}

	if len(check) > 0 && check[0] != nil {
		return check[0](r.e)
	} else {
		return r.e
	}
}

func (r Result[T]) IsErr() bool {
	return !errors.IsNil(r.e)
}

func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsErr() {
		return r
	}

	r.v = generic.Ptr(f(generic.DePtr(r.v)))
	return r
}

func (r Result[T]) Expect(msg string, args ...interface{}) T {
	if r.IsErr() {
		panic(errors.Wrapf(r.e, msg, args...))
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) OrElse(v T) T {
	if r.IsErr() {
		return v
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) OnErr(check func(err error) error) Result[T] {
	if r.IsErr() {
		r.e = check(r.e)
	}

	return r
}

func (r Result[T]) Unwrap(check ...func(err error) error) T {
	if !r.IsErr() {
		return generic.DePtr(r.v)
	}

	if len(check) > 0 && check[0] != nil {
		panic(check[0](r.e))
	} else {
		panic(r.e)
	}
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", r.v)
	}

	return r.e.Error()
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, r.e
	}

	return json.Marshal(generic.DePtr(r.v))
}
