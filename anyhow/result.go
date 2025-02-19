package anyhow

import (
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/anyhow/aherrcheck"
	"github.com/pubgo/funk/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Result[T any] struct {
	v   *T
	Err Error
}

func (r Result[T]) GetValue() T {
	return r.getValue()
}

func (r Result[T]) OnValue(fn func(v T)) Result[T] {
	if !r.IsErr() {
		fn(r.getValue())
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

func (r Result[T]) Unwrap(setter *Error, callbacks ...func(err error) error) T {
	if setter == nil {
		debug.PrintStack()
		panic("Unwrap: setter is nil")
	}

	var ret = r.getValue()
	if !r.IsErr() {
		return ret
	}

	// err No checking, repeat setting
	if (*setter).IsErr() {
		log.Warn().Msgf("Unwrap: setter is not nil, err=%v", (*setter).getErr())
	}

	callbacks = append(callbacks, aherrcheck.GetErrChecks()...)
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
	if !r.IsErr() {
		return nil
	}
	
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
