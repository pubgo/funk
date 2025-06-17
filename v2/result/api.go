package result

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/errors"
)

func All[T any](results ...Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, result := range results {
		if result.IsErr() {
			return Fail[[]T](result.GetErr())
		}
		values = append(values, result.getValue())
	}
	return OK(values)
}

func Recovery(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		errMust(errors.Errorf("setter is nil"))
		return
	}

	*setter = errRecovery(
		func() bool { return *setter != nil },
		func() error { return *setter },
		callbacks...,
	)
}

func RecoveryErr(setter ErrSetter, callbacks ...func(err error) error) {
	if setter == nil {
		errMust(errors.Errorf("setter is nil"))
		return
	}

	setter.setError(errRecovery(
		func() bool { return setter.IsErr() },
		func() error { return setter.GetErr() },
		callbacks...,
	))
}

func ErrorOf(msg string, args ...any) Error {
	return newError(errors.WrapCaller(fmt.Errorf(msg, args...), 1))
}

func ErrProxyOf(err *error) *ErrProxy {
	if err == nil {
		errMust(errors.Errorf("setter is nil"))
		return &ErrProxy{}
	}
	return &ErrProxy{err: err}
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
	if err == nil {
		return Error{}
	}

	err = errors.WrapCaller(err, 1)
	return newError(err)
}

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v}
}

func Fail[T any](err error) Result[T] {
	if err == nil {
		return Result[T]{}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func Wrap[T any](v T, err error) Result[T] {
	if err == nil {
		return Result[T]{v: &v}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func WrapFn[T any](fn func() (T, error)) Result[T] {
	v, err := try1(fn)
	if err == nil {
		return Result[T]{v: &v}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func CatchErr(setter ErrSetter, err error, contexts ...context.Context) bool {
	return catchErr(newError(err), setter, nil, contexts...)
}

func Catch(rawSetter *error, err error, contexts ...context.Context) bool {
	return catchErr(newError(err), nil, rawSetter, contexts...)
}

func Try[T any](fn func() (T, error)) (r Result[T]) {
	if fn == nil {
		err := errors.WrapCaller(errors.New("function is nil"), 1)
		return Fail[T](err)
	}

	defer Recovery(&r.err)
	return Wrap(fn())
}

func MapTo[T, U any](r Result[T], fn func(T) U) Result[U] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Fail[U](err)
	}

	return OK(fn(r.getValue()))
}

func FlatMapTo[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Fail[U](err)
	}

	return fn(r.getValue())
}
