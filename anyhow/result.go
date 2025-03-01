package anyhow

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/samber/lo"
)

type Result[T any] struct {
	_ [0]func() // disallow ==

	v   *T
	Err Error
}

func (r Result[T]) GetValue() T {
	if r.IsErr() {
		debug.PrintStack()
		err := errors.WrapCaller(r.getErr(), 1)
		errors.Debug(err)
		panic(err)
	}
	return r.getValue()
}

func (r Result[T]) OnValue(fn func(v T)) Result[T] {
	if !r.IsErr() {
		fn(r.getValue())
	}
	return r
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
		var err = errors.WrapCaller(r.getErr(), 1)
		return err
	}

	if v == nil {
		debug.PrintStack()
		panic("v params is nil")
	}

	*v = r.getValue()
	return nil
}

func (r Result[T]) Expect(format string, args ...any) T {
	if !r.IsErr() {
		return r.getValue()
	}

	debug.PrintStack()
	err := errors.WrapCaller(r.getErr(), 1)
	err = errors.Wrapf(err, format, args...)
	errors.Debug(err)
	panic(err)
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", r.getValue())
	}

	return fmt.Sprint(errors.WrapCaller(r.getErr(), 1))
}

func (r Result[T]) RawErrTo(setter *error, ctx ...context.Context) bool {
	return errTo(r.Err, nil, setter, ctx...)
}

func (r Result[T]) ErrTo(setter *Error, ctx ...context.Context) bool {
	return errTo(r.Err, setter, nil, ctx...)
}

func (r Result[T]) OrElse(t T) T {
	if r.IsErr() {
		return t
	}

	return r.getValue()
}

func (r Result[T]) IsErrNil() bool {
	return r.getErr() == nil
}

func (r Result[T]) IsErr() bool {
	return r.getErr() != nil
}

func (r Result[T]) GetErr() error {
	if !r.IsErr() {
		return nil
	}

	var err = r.getErr()
	return errors.WrapCaller(err, 1)
}

func (r Result[T]) WithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{Err: newError(err)}
}

func (r Result[T]) OnErr(callbacks ...func(err error) error) Result[T] {
	if r.IsErr() {
		var err = r.getErr()
		for _, fn := range callbacks {
			err = fn(err)
			if err == nil {
				return OK(r.getValue())
			}
		}
		return Wrap(r.getValue(), errors.WrapCaller(err, 1))
	}

	return r
}

func (r Result[T]) getValue() T {
	return lo.FromPtr(r.v)
}

func (r Result[T]) getErr() error {
	return r.Err.getErr()
}
