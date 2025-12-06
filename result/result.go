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
	_ Checkable = new(Result[any])
	_ ErrSetter = new(Result[any])
)

type Result[T any] struct {
	_ [0]func() // disallow ==

	v   *T
	err error
}

func (r Result[T]) WithFn(fn func() (T, error)) Result[T] {
	if r.IsErr() {
		return Result[T]{err: errors.WrapCaller(r.getErr(), 1)}
	}

	return WrapFn(fn)
}

func (r Result[T]) WithValue(v T) Result[T] {
	if r.IsErr() {
		return Result[T]{err: errors.WrapCaller(r.getErr(), 1)}
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
	return Error{}
}

func (r Result[T]) UnwrapOrLog(events ...func(e *zerolog.Event)) T {
	if r.IsErr() {
		panicIfError(errors.WrapCaller(r.getErr(), 1), events...)
	}

	return r.getValue()
}

// TryUnwrap attempts to unwrap the value, returning it and a boolean indicating success
// This method is useful when you want to safely extract a value from a Result
// without triggering a panic. It returns the value (or zero value if error)
// and a boolean indicating whether the extraction was successful.
//
// Example:
//
//	if value, ok := result.TryUnwrap(); ok {
//	    fmt.Printf("Success: %v\n", value)
//	} else {
//	    fmt.Println("Operation failed")
//	}
func (r Result[T]) TryUnwrap() (T, bool) {
	if r.IsErr() {
		var zero T
		return zero, false
	}
	return r.getValue(), true
}

// Match allows pattern matching on the Result, applying the appropriate function
// This method provides a way to handle both success and error cases in a single operation
// similar to pattern matching in functional languages.
//
// Example:
//
//	result.OK(42).Match(
//	    func(value int) { fmt.Printf("Success: %d\n", value) },
//	    func(err error) { fmt.Printf("Error: %v\n", err) },
//	)
func (r Result[T]) Match(onOk func(T), onErr func(error)) {
	if r.IsOK() {
		onOk(r.getValue())
	} else {
		onErr(r.getErr())
	}
}

// MatchWithResult allows pattern matching with a result-returning function
// This method is similar to Match, but both handler functions return a Result[T],
// allowing for chaining operations that may themselves produce Results.
//
// Example:
//
//	result.OK(42).MatchWithResult(
//	    func(value int) result.Result[int] { return result.OK(value * 2) },
//	    func(err error) result.Result[int] { return result.Fail[int](err) },
//	)
func (r Result[T]) MatchWithResult(onOk func(T) Result[T], onErr func(error) Result[T]) Result[T] {
	if r.IsOK() {
		return onOk(r.getValue())
	}
	return onErr(r.getErr())
}

func (r Result[T]) Must(events ...func(e *zerolog.Event)) {
	if r.IsErr() {
		panicIfError(errors.WrapCaller(r.getErr(), 1), events...)
	}
}

func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panicIfError(errors.WrapCaller(r.getErr(), 1))
	}
	return r.getValue()
}

func (r Result[T]) UnwrapErr() (T, error) {
	if r.IsErr() {
		var zero T
		return zero, r.getErr()
	}
	return r.getValue(), nil
}

func (r Result[T]) Or(defaultVal T) Result[T] {
	if r.IsErr() {
		return OK(defaultVal)
	}
	return r
}

func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsErr() {
		return defaultVal
	}
	return r.getValue()
}

func (r Result[T]) OrElse(fn func() T) Result[T] {
	if r.IsErr() {
		return OK(fn())
	}
	return r
}

func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.IsErr() {
		return fn()
	}
	return r.getValue()
}

func (r Result[T]) UnwrapOrEmpty() (t T) {
	if r.IsErr() {
		return
	}
	return r.getValue()
}

func (r Result[T]) ThrowErr(setter ErrSetter, contexts ...context.Context) bool {
	return catchErr(ErrOf(r.err), setter, nil, contexts...)
}

func (r Result[T]) UnwrapOrThrow(setter ErrSetter, contexts ...context.Context) (t T) {
	ret, err := unwrapErr(r, nil, setter, contexts...)
	if err != nil {
		setError(setter, errors.WrapCaller(err, 1))
	}
	return ret
}

func (r Result[T]) CallIfOK(fn func(val T) error) Result[T] {
	if r.IsOK() {
		return Fail[T](fn(r.getValue()))
	}
	return r
}

func (r Result[T]) Expect(format string, args ...any) T {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		panicIfError(err, func(e *zerolog.Event) {
			e.Str(logfields.Msg, fmt.Sprintf(format, args...))
		})
	}

	return r.getValue()
}

func (r Result[T]) IsErr() bool { return r.getErr() != nil }

func (r Result[T]) IsOK() bool { return r.getErr() == nil }

func (r Result[T]) InspectErr(fn func(err error)) {
	if r.IsErr() {
		fn(r.getErr())
	}
}

func (r Result[T]) Inspect(fn func(val T)) {
	if r.IsOK() {
		fn(r.getValue())
	}
}

// IfErr executes fn if the result is an error, then returns the result unchanged.
// This is similar to InspectErr but allows chaining with other operations.
func (r Result[T]) IfErr(fn func(err error)) Result[T] {
	if r.IsErr() {
		fn(r.getErr())
	}
	return r
}

// IfOK executes fn if the result is OK, then returns the result unchanged.
// This is similar to Inspect but allows chaining with other operations.
func (r Result[T]) IfOK(fn func(val T)) Result[T] {
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

func (r Result[T]) GetErr() error { return r.Err() }

func (r Result[T]) Err() error {
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

func (r Result[T]) WithErr(err error, tags ...errors.Tags) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapTagsCaller(err, lo.FirstOrEmpty(tags), 1)
	return Result[T]{err: err}
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
