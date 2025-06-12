package result

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
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
	return Result[T]{E: errors.WrapCaller(err, 1)}
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, E: errors.WrapCaller(err, 1)}
}

func Of[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, E: errors.WrapCaller(err, 1)}
}

type Result[T any] struct {
	v *T
	E error
}

func (r Result[T]) WithErr(err error) Result[T] {
	return Result[T]{E: errors.WrapCaller(err, 1)}
}

func (r Result[T]) WithVal(v T) Result[T] {
	return OK(v)
}

func (r Result[T]) ValueTo(v *T) error {
	if r.IsErr() {
		return errors.WrapCaller(r.E, 1)
	}

	*v = generic.FromPtr(r.v)
	return nil
}

func (r Result[T]) OnValue(fn func(t T) error) error {
	if r.IsErr() {
		return r.E
	}

	return errors.WrapCaller(fn(generic.FromPtr(r.v)), 1)
}

func (r Result[T]) OnErr(check func(err error)) {
	if r.IsOK() {
		return
	}

	check(r.E)
}

func (r Result[T]) GetErr() error {
	if r.IsOK() {
		return nil
	}

	return errors.WrapCaller(r.E, 1)
}

func (r Result[T]) Err(check ...func(err error) error) error {
	if r.IsOK() {
		return nil
	}

	if len(check) > 0 && check[0] != nil {
		return errors.WrapCaller(check[0](r.E), 1)
	}

	return errors.WrapCaller(r.E, 1)
}

func (r Result[T]) IsErr() bool { return r.E != nil }

func (r Result[T]) IsOK() bool { return r.E == nil }

func (r Result[T]) OrElse(v T) T {
	if r.IsErr() {
		return v
	}
	return generic.FromPtr(r.v)
}

func (r Result[T]) UnwrapErr(setter *error) T {
	if setter == nil {
		debug.PrintStack()
		panic("UnwrapErr: setter is nil")
	}

	if r.IsErr() {
		*setter = errors.WrapCaller(r.E, 1)
	}

	return lo.FromPtr(r.v)
}

func (r Result[T]) Unwrap(check ...func(err error) error) T {
	if r.IsOK() {
		return generic.FromPtr(r.v)
	}

	if len(check) > 0 && check[0] != nil {
		panic(check[0](r.E))
	} else {
		panic(r.E)
	}
}

func (r Result[T]) GetValue() T {
	if r.IsOK() {
		return generic.FromPtr(r.v)
	}

	panic(errors.WrapStack(r.E))
}

func (r Result[T]) Expect(format string, args ...any) T {
	if r.IsOK() {
		return generic.FromPtr(r.v)
	}

	panic(errors.WrapStack(errors.Wrapf(r.E, format, args...)))
}

func (r Result[T]) String() string {
	if r.IsOK() {
		return fmt.Sprintf("%v", generic.FromPtr(r.v))
	}

	return fmt.Sprint(errors.WrapCaller(r.E, 1))
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, errors.WrapCaller(r.E, 1)
	}

	return json.Marshal(generic.FromPtr(r.v))
}

func (r Result[T]) UnmarshalJSON([]byte) error {
	panic("unimplemented")
}

func (r Result[T]) Do(fn func(v T)) {
	if r.IsErr() {
		return
	}

	fn(generic.FromPtr(r.v))
}

func (r Result[T]) CatchTo(setter *error, callbacks ...func(err error) error) bool {
	if setter == nil {
		debug.PrintStack()
		panic("CatchTo: setter is nil")
	}

	if r.IsOK() {
		return false
	}

	// setter err is not nil
	if *setter != nil {
		log.Err(*setter).Msgf("CatchTo: setter error is not nil")
		return true
	}

	var err = r.E
	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return false
		}
	}

	*setter = errors.WrapCaller(err, 1)
	return true
}

func (r Result[T]) InspectErr(fn func(error)) Result[T] {
	if r.IsErr() {
		fn(r.E)
	}
	return r
}

func (r Result[T]) Inspect(fn func(T)) Result[T] {
	if r.IsOK() {
		fn(generic.FromPtr(r.v))
	}
	return r
}

func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	if r.IsOK() {
		return r
	}
	return fn(generic.FromPtr(r.v))
}

func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsOK() {
		return r
	}
	return OK(fn(generic.FromPtr(r.v)))
}

func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.IsOK() {
		return r
	}
	return Err[T](fn(r.E))
}

func MapTo[Src any, To any](s Result[Src], do func(s Src) To) Result[To] {
	if s.IsErr() {
		return Err[To](errors.WrapCaller(s.Err(), 1))
	}

	return OK(do(s.Unwrap()))
}

func FlatMap[Src any, To any](s Result[Src], do func(s Src) (r Result[To])) Result[To] {
	if s.IsErr() {
		return Err[To](errors.WrapCaller(s.Err(), 1))
	}

	return do(s.Unwrap())
}

func Unwrap[T any](ret Result[T], gErr *error, callback ...func(err error) error) T {
	if gErr == nil {
		debug.PrintStack()
		panic("Unwrap: gErr is nil")
	}

	if ret.IsOK() {
		return ret.GetValue()
	}

	var t T
	err := ret.Err()
	for _, fn := range callback {
		if err == nil {
			return t
		}

		err = fn(err)
	}

	*gErr = errors.WrapCaller(err, 1)
	return t
}
