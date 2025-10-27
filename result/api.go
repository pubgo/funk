package result

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/v2/errors"
)

func Run(executors ...func() error) Error {
	for _, executor := range executors {
		if err := executor(); err != nil {
			return ErrOf(errors.WrapCaller(err, 1))
		}
	}
	return Error{}
}

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

func RecoveryErr(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		errNilOrPanic(errors.Errorf("setter is nil"))
		return
	}

	setError(ErrProxyOf(setter), errRecovery(
		func() error { return *setter },
		callbacks...,
	))
}

func Recovery(setter ErrSetter, callbacks ...func(err error) error) {
	if setter == nil {
		errNilOrPanic(errors.Errorf("setter is nil"))
		return
	}

	setError(setter, errRecovery(
		func() error { return setter.GetErr() },
		callbacks...,
	))
}

func Errorf(msg string, args ...any) Error {
	return newError(errors.WrapCaller(fmt.Errorf(msg, args...), 1))
}

func ErrProxyOf(err *error) ErrProxy {
	if err == nil {
		errNilOrPanic(errors.Errorf("err param is nil"))
		return ErrProxy{}
	}
	return ErrProxy{err: err}
}

func ErrOf(err error) Error {
	if err == nil {
		return Error{}
	}

	err = errors.WrapCaller(err, 1)
	return newError(err)
}

func ErrOfFn(fn func() error) Error {
	err := try(fn)
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

func Catch(setter ErrSetter, err error, contexts ...context.Context) bool {
	return catchErr(newError(err), setter, nil, contexts...)
}

func CatchErr(rawSetter *error, err error, contexts ...context.Context) bool {
	return catchErr(newError(err), nil, rawSetter, contexts...)
}

func MapTo[T, U any](r Result[T], fn func(T) U) Result[U] {
	if r.IsErr() {
		return Fail[U](errors.WrapCaller(r.getErr(), 1))
	}

	return OK(fn(r.getValue()))
}

func FlatMapTo[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return Fail[U](errors.WrapCaller(r.getErr(), 1))
	}

	return fn(r.getValue())
}

func LogErr(err error, events ...func(e *zerolog.Event)) {
	logErr(nil, 0, err, events...)
}

func LogErrCtx(ctx context.Context, err error, events ...func(e *zerolog.Event)) {
	logErr(ctx, 0, err, events...)
}

func Must(err error, events ...func(e *zerolog.Event)) {
	if err == nil {
		return
	}

	errNilOrPanic(errors.WrapCaller(err, 1), events...)
}

func Must1[T any](ret T, err error) T {
	if err != nil {
		errNilOrPanic(errors.WrapCaller(err, 1))
	}

	return ret
}
