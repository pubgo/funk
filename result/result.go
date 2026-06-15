package result

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/errors"
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

// WrapFn wraps a function into a Result.
func (r Result[T]) WithFn(fn func() (T, error)) Result[T] {
	if r.IsErr() {
		return Result[T]{err: errors.WrapCaller(r.getErr(), 1)}
	}

	return WrapFn(fn)
}

// WithValue wraps a value into a Result.
func (r Result[T]) WithValue(v T) Result[T] {
	if r.IsErr() {
		return Result[T]{err: errors.WrapCaller(r.getErr(), 1)}
	}

	return OK(v)
}

// ValueTo extracts the value from the Result and assigns it to the provided variable.
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

// UnwrapOrLog attempts to unwrap the value, returning it and panicking if an error occurs.
func (r Result[T]) UnwrapOrLog(events ...func(e Event)) T {
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

// Must panics if the result is an error.
func (r Result[T]) Must(events ...func(e Event)) {
	if r.IsErr() {
		panicIfError(errors.WrapCaller(r.getErr(), 1), events...)
	}
}

// Unwrap returns the value if it's OK, or panics if it's an error.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panicIfError(errors.WrapCaller(r.getErr(), 1))
	}
	return r.getValue()
}

// UnwrapErr returns the value and a nil error when the result is OK,
// or the zero value and the underlying error when the result failed.
func (r Result[T]) UnwrapErr() (T, error) {
	if r.IsErr() {
		var zero T
		return zero, r.getErr()
	}
	return r.getValue(), nil
}

// Or returns a new Result with the provided default value if the current Result is an error.
func (r Result[T]) Or(defaultVal T) Result[T] {
	if r.IsErr() {
		return OK(defaultVal)
	}
	return r
}

// UnwrapOr returns the value if it's OK, or the default value if it's an error.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsErr() {
		return defaultVal
	}
	return r.getValue()
}

// OrElse returns a new Result with the provided default value if the current Result is an error.
func (r Result[T]) OrElse(fn func() T) Result[T] {
	if r.IsErr() {
		return OK(fn())
	}
	return r
}

// UnwrapOrElse returns the value if it's OK, or the result of the provided function if it's an error.
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.IsErr() {
		return fn()
	}
	return r.getValue()
}

// UnwrapOrEmpty returns the value if it's OK, or the zero value if it's an error.
func (r Result[T]) UnwrapOrEmpty() (t T) {
	if r.IsErr() {
		return t
	}
	return r.getValue()
}

// ThrowErr throws an error if the result is an error.
func (r Result[T]) ThrowErr(err *error, contexts ...context.Context) bool {
	return catchErr(ErrOf(r.err), nil, err, contexts...)
}

// Throw returns the value if it's OK, or throws an error if it's an error.
func (r Result[T]) Throw(setter ErrSetter, contexts ...context.Context) bool {
	return catchErr(ErrOf(r.err), setter, nil, contexts...)
}

// UnwrapOrThrow returns the value if it's OK, or throws an error if it's an error.
func (r Result[T]) UnwrapOrThrow(setter ErrSetter, contexts ...context.Context) (t T) {
	ret, err := unwrapErr(r, nil, setter, contexts...)
	if err != nil {
		setError(setter, errors.WrapCaller(err, 1))
	}
	return ret
}

// CallIfOK calls fn with the value when the result is OK and propagates any returned error.
func (r Result[T]) CallIfOK(fn func(val T) error) Result[T] {
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

// Expect panics if the result is an error.
func (r Result[T]) Expect(format string, args ...any) T {
	if r.IsErr() {
		err := errors.WrapCaller(r.getErr(), 1)
		panicIfError(err, func(e Event) {
			e.Msgf(format, args...)
		})
	}

	return r.getValue()
}

// IsErr returns true if the result is an error.
func (r Result[T]) IsErr() bool { return r.getErr() != nil }

// IsOK returns true if the result is OK.
func (r Result[T]) IsOK() bool { return r.getErr() == nil }

// InspectErr executes fn if the result is an error.
func (r Result[T]) InspectErr(fn func(err error)) {
	if r.IsErr() {
		fn(r.getErr())
	}
}

// Inspect executes fn if the result is OK, then returns the result unchanged.
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

// LogCtx logs the error with the provided context.
func (r Result[T]) LogCtx(ctx context.Context, events ...func(e Event)) Result[T] {
	logErr(ctx, 0, r.err, events...)
	return r
}

// Log logs the error.
func (r Result[T]) Log(events ...func(e Event)) Result[T] {
	logErr(context.Background(), 0, r.err, events...)
	return r
}

// Validate runs fn on the successful value and returns a failed result when validation fails.
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

// Map transforms the successful value and preserves any existing error.
func (r Result[T]) Map(fn func(val T) T) Result[T] {
	if r.IsErr() {
		return r
	}
	return OK(fn(r.getValue()))
}

// MapVal transforms the successful value with fn and propagates the first error.
func (r Result[T]) MapVal(fn func(val T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return fn(r.getValue())
}

// FlatMap is an alias of MapVal for callers who prefer the FlatMap naming convention.
func (r Result[T]) FlatMap(fn func(val T) Result[T]) Result[T] {
	return r.MapVal(fn)
}

// MapErr transforms the error when the result failed and leaves successful values unchanged.
func (r Result[T]) MapErr(fn func(err error) error) Result[T] {
	if r.IsOK() {
		return r
	}
	return Fail[T](fn(r.getErr()))
}

// MapErrOr transforms the error with fn when the result failed and returns the produced result.
func (r Result[T]) MapErrOr(fn func(err error) Result[T]) Result[T] {
	if r.IsOK() {
		return r
	}
	return fn(r.getErr())
}

// GetErr returns the error if the result is an error, or nil if it's OK.
func (r Result[T]) GetErr() error { return r.Err() }

// Err returns the error if the result is an error, or nil if it's OK.
func (r Result[T]) Err() error {
	if r.IsOK() {
		return nil
	}

	return r.getErr()
}

// String returns a string representation of the result.
func (r Result[T]) String() string {
	if r.IsOK() {
		return fmt.Sprintf("OK(%v)", r.getValue())
	}
	return fmt.Sprintf("Error(%v)", r.getErr())
}

// WithErrorf replaces the result with a formatted error and discards any successful value.
func (r Result[T]) WithErrorf(format string, args ...any) Result[T] {
	err := fmt.Errorf(format, args...)
	err = errors.WrapCaller(err, 1)
	return Result[T]{err: err}
}

// WithErr returns a new Result with the provided error.
func (r Result[T]) WithErr(err error, tags ...errors.Tags) Result[T] {
	if err == nil {
		return r
	}

	err = errors.WrapTagsCaller(err, lo.FirstOrEmpty(tags), 1)
	return Result[T]{err: err}
}

// MarshalJSON returns a JSON representation of the result.
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

func (r *Result[T]) applyErr(err error) {
	if r == nil || err == nil {
		return
	}
	r.err = err
	r.v = nil
}
