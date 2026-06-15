package result

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
	"github.com/pubgo/funk/v2/stack"
)

func Run(executors ...func() error) Error {
	for _, executor := range executors {
		if err := executor(); err != nil {
			return newError(errors.WrapCaller(err, 1))
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
		panicIfError(errors.Errorf("setter is nil"))
		return
	}

	err := errparser.Parse(recover())
	if err == nil {
		err = *setter
	}

	if err == nil {
		return
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	stack.Print()
	setError(ErrProxyOf(setter), errors.WrapCaller(err, 1))
}

func Recovery(setter ErrSetter, callbacks ...func(err error) error) {
	if setter == nil {
		panicIfError(errors.Errorf("setter is nil"))
		return
	}

	err := errparser.Parse(recover())
	if err == nil {
		err = setter.GetErr()
	}

	if err == nil {
		return
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	stack.Print()
	setError(setter, errors.WrapCaller(err, 1))
}

func Errorf(msg string, args ...any) Error {
	return newError(errors.WrapCaller(fmt.Errorf(msg, args...), 1))
}

func ProxyOf(err *error) ProxyErr { return ErrProxyOf(err) }
func ErrProxyOf(err *error) ProxyErr {
	if err == nil {
		panicIfError(errors.Errorf("err param is nil"))
		return ProxyErr{}
	}
	return ProxyErr{err: err}
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
		panicIfError(errors.WrapCaller(errors.New("result.Fail called with nil error"), 1))
	}

	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

func WrapErr[T any](v T, err error) (t T, gErr Error) {
	if err == nil {
		return v, gErr
	}

	return t, newError(errors.WrapCaller(err, 1))
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

func Throw(setter ErrSetter, err error, contexts ...context.Context) bool {
	err = errors.WrapCaller(err, 1)
	return catchErr(newError(err), setter, nil, contexts...)
}

func ThrowErr(rawSetter *error, err error, contexts ...context.Context) bool {
	err = errors.WrapCaller(err, 1)
	return catchErr(newError(err), nil, rawSetter, contexts...)
}

func MapTo[T, U any](r Result[T], fn func(T) U) Result[U] {
	if r.IsErr() {
		return Fail[U](errors.WrapCaller(r.getErr(), 1))
	}

	return OK(fn(r.getValue()))
}

func MapValTo[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return Fail[U](errors.WrapCaller(r.getErr(), 1))
	}

	return fn(r.getValue())
}

// FlatMapTo transforms a successful value with fn and propagates the first error.
// It is an alias of MapValTo for callers who prefer the FlatMap naming convention.
func FlatMapTo[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	return MapValTo(r, fn)
}

func LogErr(err error, events ...func(e Event)) {
	err = errors.WrapCaller(err, 1)
	logErr(context.Background(), 0, err, events...)
}

func LogErrCtx(ctx context.Context, err error, events ...func(e Event)) {
	err = errors.WrapCaller(err, 1)
	logErr(ctx, 0, err, events...)
}

func Must(err error, events ...func(e Event)) {
	if err == nil {
		return
	}

	panicIfError(errors.WrapCaller(err, 1), events...)
}

func Must1[T any](ret T, err error) T {
	if err != nil {
		panicIfError(errors.WrapCaller(err, 1))
	}

	return ret
}

// Try converts a function that may panic into a Result
// This function wraps a potentially panicking function and converts panics
// into error Results, providing a safer way to call functions that might panic.
//
// Example:
//
//	result := result.Try(func() int {
//	    // Some operation that might panic
//	    return riskyOperation()
//	})
//	if value, ok := result.TryUnwrap(); ok {
//	    fmt.Printf("Success: %d\n", value)
//	}
func Try[T any](fn func() T) (r Result[T]) {
	defer func() {
		if err := recover(); err != nil {
			r = Fail[T](fmt.Errorf("panic occurred: %w", errparser.Parse(err)))
		}
	}()
	return OK(fn())
}

// Partition separates a slice of Results into values and errors
// This function takes a slice of Results and separates them into two slices:
// one containing all successful values and another containing all errors.
// This is useful when you want to process all successes and handle all errors
// separately rather than stopping at the first error.
//
// Example:
//
//	results := []result.Result[int]{result.OK(1), result.Fail[int](err1), result.OK(2)}
//	values, errors := result.Partition(results)
//	fmt.Printf("Values: %v, Errors: %v\n", values, errors)
func Partition[T any](results []Result[T]) ([]T, []error) {
	var values []T
	var errs []error

	for _, result := range results {
		if result.IsOK() {
			values = append(values, result.getValue())
		} else {
			errs = append(errs, result.GetErr())
		}
	}

	return values, errs
}

// Collect takes a slice of Results and returns either all values or the first error
// This function processes a slice of Results and collects all successful values
// into a single Result containing a slice of values. If any Result contains an error,
// the function stops and returns that error in a Result.
//
// Example:
//
//	results := []result.Result[int]{result.OK(1), result.OK(2), result.OK(3)}
//	collected := result.Collect(results)
//	if vals, ok := collected.TryUnwrap(); ok {
//	    fmt.Printf("Collected values: %v\n", vals)
//	}
func Collect[T any](results []Result[T]) Result[[]T] {
	var values = make([]T, 0, len(results))

	for _, result := range results {
		if result.IsErr() {
			return Fail[[]T](result.GetErr())
		}
		values = append(values, result.getValue())
	}

	return OK(values)
}
