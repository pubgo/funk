package anyhow

import (
	"runtime/debug"

	"github.com/pubgo/funk/errors"
)

func RecoveryErr(setter *error, callbacks ...func(err error) error) {
	recovery(
		setter,
		func() bool { return *setter != nil },
		func() error { return *setter },
		func(err error) error { return err },
		callbacks...,
	)
}

func Recovery(setter *Error, callbacks ...func(err error) error) {
	recovery(
		setter,
		func() bool { return setter.IsErr() },
		func() error { return setter.GetErr() },
		func(err error) Error { return newError(err) },
		callbacks...,
	)
}

func recovery[T any](setter *T, isErr func() bool, getErr func() error, newErr func(err error) T, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	err := errors.Parse(recover())
	if err == nil && !isErr() {
		return
	}

	if err == nil {
		err = getErr()
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	err = errors.WrapCaller(err, 1)
	*setter = newErr(err)
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

func Err[T any](err error) Result[T] {
	if err == nil {
		return Result[T]{}
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
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

func DoResult[T any](fn func() (r Result[T])) (t T, gErr error) {
	return t, try(func() error { return fn().ValueTo(&t) })
}

func DoError(fn func() (r Error)) error {
	return try(func() error { return fn().GetErr() })
}
