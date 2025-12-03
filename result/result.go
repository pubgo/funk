package result

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log/logfields"
)

var (
	_ Catchable = new(Result[any])
	_ Checkable = new(Result[any])
	_ ErrSetter = new(Result[any])
)

type Result[T any] struct {
	_ [0]func() // disallow ==

	v   *T
	err error
}

func (r Result[T]) GetValue() (t T) {
	if r.IsErr() {
		return t
	}

	return r.getValue()
}

func (r Result[T]) WithFn(fn func() (T, error)) Result[T] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Result[T]{err: err}
	}

	return WrapFn(fn)
}

func (r Result[T]) WithValue(v T) Result[T] {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		return Result[T]{err: err}
	}

	return OK(v)
}

func (r Result[T]) ValueTo(v *T) Error {
	if r.IsErr() {
		return newError(errors.WrapCaller(r.getErr(), 1))
	}

	if v == nil {
		return newError(errors.WrapStack(errors.New("v param is nil")))
	}

	*v = r.getValue()
	return newError(nil)
}

func (r Result[T]) OrElse(v T) T {
	if r.IsErr() {
		return v
	}
	return r.getValue()
}

// UnwrapOr returns the value if OK, otherwise returns the default value.
// This is a simpler alternative to Unwrap that doesn't require a setter.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsErr() {
		return defaultVal
	}
	return r.getValue()
}

// UnwrapOrElse returns the value if OK, otherwise returns the result of calling fn.
// This is useful when the default value needs to be computed lazily.
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.IsErr() {
		return fn()
	}
	return r.getValue()
}

func (r Result[T]) Expect(format string, args ...any) T {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		errNilOrPanic(err, func(e *zerolog.Event) {
			e.Str(logfields.Msg, fmt.Sprintf(format, args...))
		})
	}

	return r.getValue()
}

func (r Result[T]) Must(events ...func(e *zerolog.Event)) T {
	if r.IsErr() {
		errNilOrPanic(errors.WrapCaller(r.getErr(), 1), events...)
	}

	return r.getValue()
}

func (r Result[T]) CatchErr(setter *error, ctx ...context.Context) bool {
	return catchErr(newError(r.err), nil, setter, ctx...)
}

func (r Result[T]) Catch(setter ErrSetter, ctx ...context.Context) bool {
	return catchErr(newError(r.err), setter, nil, ctx...)
}

func (r Result[T]) IsErr() bool { return r.getErr() != nil }

func (r Result[T]) IsOK() bool { return r.getErr() == nil }

func (r Result[T]) InspectErr(fn func(err error)) Result[T] {
	if r.IsErr() {
		fn(r.getErr())
	}
	return r
}

func (r Result[T]) Inspect(fn func(val T)) Result[T] {
	if r.IsOK() {
		fn(r.getValue())
	}
	return r
}

// IfErr executes fn if the result is an error, then returns the result unchanged.
// This is similar to InspectErr but allows chaining with other operations.
func (r Result[T]) IfErr(fn func(error)) Result[T] {
	if r.IsErr() {
		fn(r.getErr())
	}
	return r
}

// IfOK executes fn if the result is OK, then returns the result unchanged.
// This is similar to Inspect but allows chaining with other operations.
func (r Result[T]) IfOK(fn func(T)) Result[T] {
	if r.IsOK() {
		fn(r.getValue())
	}
	return r
}

func (r Result[T]) LogCtx(ctx context.Context, events ...func(e *zerolog.Event)) Result[T] {
	logErr(ctx, 0, r.err, events...)
	return r
}

func (r Result[T]) Log(events ...func(e *zerolog.Event)) Result[T] {
	logErr(context.Background(), 0, r.err, events...)
	return r
}

func (r Result[T]) Map(fn func(val T) T) Result[T] {
	if r.IsErr() {
		return r
	}
	return OK(fn(r.getValue()))
}

func (r Result[T]) FlatMap(fn func(val T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return fn(r.getValue())
}

func (r Result[T]) Validate(fn func(val T) error) Result[T] {
	if r.IsErr() {
		return r
	}

	val := r.getValue()
	err := fn(val)
	if err != nil {
		return Fail[T](errors.WrapCaller(err, 1))
	}
	return OK(val)
}

func (r Result[T]) MapErr(fn func(err error) error) Result[T] {
	if r.IsOK() {
		return r
	}
	return Fail[T](fn(r.getErr()))
}

func (r Result[T]) MapErrOr(fn func(err error) Result[T]) Result[T] {
	if r.IsOK() {
		return r
	}
	return fn(r.getErr())
}

func (r Result[T]) GetErr() error {
	if r.IsOK() {
		return nil
	}

	return r.getErr()
}

func (r Result[T]) String() string {
	if r.IsOK() {
		return fmt.Sprintf("OK(%v)", r.getValue())
	}
	return fmt.Sprintf("Error(%v)", r.getErr())
}

func (r Result[T]) WithErrorf(format string, args ...any) Result[T] {
	err := fmt.Errorf(format, args...)
	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func (r Result[T]) WithErr(err error) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func (r Result[T]) WrapErr(err *errors.Err, tags ...errors.Tag) Result[T] {
	return Result[T]{err: errors.WrapTag(errors.WrapCaller(err, 1), tags...)}
}

func (r Result[T]) UnwrapErr(setter *error, contexts ...context.Context) T {
	ret, err := unwrapErr(r, setter, nil, contexts...)
	if err != nil {
		*setter = errors.WrapCaller(err, 1)
	}
	return ret
}

func (r Result[T]) Unwrap(setter ErrSetter, contexts ...context.Context) T {
	ret, err := unwrapErr(r, nil, setter, contexts...)
	if err != nil {
		setError(setter, errors.WrapCaller(err, 1))
	}
	return ret
}

// UnwrapGo returns the value and error in the standard Go style (val, err).
// This is more idiomatic for Go developers and doesn't require a setter.
// If the result is OK, returns (value, nil).
// If the result is an error, returns (zero value, error).
func (r Result[T]) UnwrapGo() (T, error) {
	if r.IsErr() {
		var zero T
		return zero, r.getErr()
	}
	return r.getValue(), nil
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, errors.WrapCaller(r.err, 1)
	}

	return json.Marshal(funk.FromPtr(r.v))
}

func (r Result[T]) getValue() T { return lo.FromPtr(r.v) }

func (r Result[T]) getErr() error { return r.err }

func (r Result[T]) setErrorInner() {
}
