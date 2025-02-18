package result

import (
	"runtime/debug"

	"github.com/pubgo/funk/anyhow"
	"github.com/pubgo/funk/errors"
)

func OK[T any](v T) Result[T] {
	return Result[T]{Result: anyhow.OK(v)}
}

func Err[T any](err error) Result[T] {
	err = errors.WrapCaller(err, 1)
	return Result[T]{Result: anyhow.Err[T](err)}
}

func ErrOf(err error) anyhow.Error {
	err = errors.WrapCaller(err, 1)
	return anyhow.ErrOf(err)
}

func ErrOfFn(fn func() error) anyhow.Error {
	return anyhow.ErrOfFn(fn)
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{Result: anyhow.Wrap(v, err)}
}

func WrapFn[T any](fn func() (T, error)) Result[T] {
	return Result[T]{Result: anyhow.WrapFn(fn)}
}

func Of[T any](v T, err error) Result[T] {
	return Result[T]{Result: anyhow.Wrap(v, err)}
}

type Result[T any] struct {
	anyhow.Result[T]
}

func (r Result[T]) WithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Err[T](err)
}

func (r Result[T]) WithVal(v T) Result[T] {
	if r.IsErr() {
		return Err[T](errors.WrapCaller(r.Err(), 1))
	}

	return OK(v)
}

func (r Result[T]) ValueTo(v *T) error {
	if r.IsErr() {
		return errors.WrapCaller(r.Result.GetErr(), 1)
	}

	if v == nil {
		debug.PrintStack()
		panic("v params is nil")
	}

	*v = r.Result.GetValue()
	return nil
}

func (r Result[T]) OnValue(fn func(t T) error) error {
	if r.IsErr() {
		return errors.WrapCaller(r.Result.GetErr(), 1)
	}

	return errors.WrapCaller(fn(r.Result.GetValue()), 1)
}

func (r Result[T]) ErrTo(gErr *anyhow.Error, callback ...func(err error) error) T {
	if !r.Result.IsErr() {
		return r.Result.GetValue()
	}

	var err = errors.WrapCaller(r.Result.GetErr(), 1)
	for _, fn := range callback {
		err = fn(err)
		if err == nil {
			return r.Result.GetValue()
		}
	}

	*gErr = anyhow.ErrOf(err)
	return r.Result.GetValue()
}

func (r Result[T]) Unwrap(callback ...func(err error) error) T {
	if !r.Result.IsErr() {
		return r.Result.GetValue()
	}

	var err = errors.WrapCaller(r.Result.GetErr(), 1)
	for _, fn := range callback {
		err = fn(err)
		if err == nil {
			return r.Result.GetValue()
		}
	}

	debug.PrintStack()
	panic(err)
}

func (r Result[T]) Expect(format string, args ...any) T {
	return r.Result.Expect(format, args...)
}

func (r Result[T]) String() string {
	return r.Result.String()
}

func (r Result[T]) Do(fn func(v T)) Result[T] {
	return Result[T]{Result: r.Result.OnValue(fn)}
}

func (r Result[T]) IsErr() bool {
	return r.Result.IsErr()
}

func (r Result[T]) Err() error {
	return r.Result.GetErr()
}
