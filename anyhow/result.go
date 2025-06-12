package anyhow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/errors"
	"github.com/samber/lo"
)

type Result[T any] struct {
	_ [0]func() // disallow ==

	v   *T
	Err error
}

func (r Result[T]) GetValue() T {
	if r.IsErr() {
		errMust(errors.WrapCaller(r.getErr(), 1))
	}

	return r.getValue()
}

func (r Result[T]) SetWithValue(v T) Result[T] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Result[T]{Err: err}
	}

	return OK(v)
}

func (r Result[T]) ValueTo(v *T) Error {
	if r.IsErr() {
		return newError(errors.WrapCaller(r.getErr(), 1))
	}

	if v == nil {
		return newError(errors.WrapStack(errors.New("v params is nil")))
	}

	*v = r.getValue()
	return newError(nil)
}

func (r Result[T]) Expect(format string, args ...any) T {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		errMust(errors.Wrapf(err, format, args...))
	}

	return r.getValue()
}

func (r Result[T]) Must() T {
	if r.IsErr() {
		errMust(errors.WrapCaller(r.getErr(), 1))
	}

	return r.getValue()
}

func (r Result[T]) Catch(setter *error, ctx ...context.Context) bool {
	return catchErr(newError(r.Err), nil, setter, ctx...)
}

func (r Result[T]) CatchErr(setter *Error, ctx ...context.Context) bool {
	return catchErr(newError(r.Err), setter, nil, ctx...)
}

func (r Result[T]) IsErr() bool { return r.getErr() != nil }

func (r Result[T]) IsOK() bool { return r.getErr() == nil }

func (r Result[T]) Filter(predicate func(T) bool, errorMsg string) Result[T] {
	if r.IsErr() {
		return r
	}

	if predicate(r.getValue()) {
		return r
	}

	return Fail[T](errors.New(errorMsg))
}

func (r Result[T]) InspectErr(fn func(error)) Result[T] {
	if r.IsErr() {
		fn(r.getErr())
	}
	return r
}

func (r Result[T]) Inspect(fn func(T)) Result[T] {
	if r.IsOK() {
		fn(r.getValue())
	}
	return r
}

func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsOK() {
		return r
	}
	return OK(fn(r.getValue()))
}

func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.IsOK() {
		return r
	}
	return Fail[T](fn(r.getErr()))
}

func (r Result[T]) GetErr() error {
	if r.IsOK() {
		return nil
	}

	err := r.getErr()
	return errors.WrapCaller(err, 1)
}

func (r Result[T]) String() string {
	if r.IsOK() {
		return fmt.Sprintf("Ok(%v)", r.getValue())
	}
	return fmt.Sprintf("Error(%v)", r.getErr())
}

func (r Result[T]) SetWithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: err}
}

func (r Result[T]) Unwrap(setter *error, contexts ...context.Context) T {
	ret, err := unwrapErr(r, setter, nil, contexts...)
	if err != nil {
		*setter = errors.WrapCaller(err, 1)
	}
	return ret
}

func (r Result[T]) UnwrapErr(setter *Error, contexts ...context.Context) T {
	ret, err := unwrapErr(r, nil, setter, contexts...)
	if err != nil {
		*setter = newError(errors.WrapCaller(err, 1))
	}
	return ret
}

func (r Result[T]) OrElse(fn func(error) T) Result[T] {
	if r.IsOK() {
		return r
	}
	return OK(fn(r.getErr()))
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.IsOK() {
		return r.getValue()
	}
	return defaultValue
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, errors.WrapCaller(r.Err, 1)
	}

	return json.Marshal(funk.FromPtr(r.v))
}

func (r Result[T]) getValue() T { return lo.FromPtr(r.v) }

func (r Result[T]) getErr() error { return r.Err }
