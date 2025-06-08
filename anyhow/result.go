package anyhow

import (
	"fmt"

	"github.com/pubgo/funk/errors"
)

// Result represents a value that might be an error, similar to Rust's Result<T, E>
// It implements multiple interfaces for functional programming patterns
type Result[T any] struct {
	value T
	err   error
	isOk  bool
}

// Ensure Result implements all relevant interfaces at compile time
var _ Checkable = (*Result[int])(nil)
var _ ErrorAccessible = (*Result[int])(nil)
var _ Unwrappable[int] = (*Result[int])(nil)
var _ Panicable[int] = (*Result[int])(nil)
var _ Catchable = (*Result[int])(nil)

// === Result Constructors ===

// Ok creates a successful Result containing the given value
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, isOk: true}
}

// Fail creates a failed Result containing the given error
func Fail[T any](err error) Result[T] {
	if err == nil {
		var zero T
		return Result[T]{value: zero, isOk: true}
	}
	return Result[T]{err: errors.WrapCaller(err, 1), isOk: false}
}

// From creates a Result from a value and error pair (common Go pattern)
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Fail[T](err)
	}
	return Ok(value)
}

// Try safely executes a function that returns (T, error) and wraps it in a Result
func Try[T any](fn func() (T, error)) Result[T] {
	if fn == nil {
		return Fail[T](errors.New("function is nil"))
	}

	defer func() {
		if r := recover(); r != nil {
			// Handle panic recovery if needed
		}
	}()

	value, err := fn()
	return From(value, err)
}

// === State Checking ===

// IsOk returns true if the Result contains a value
func (r Result[T]) IsOk() bool {
	return r.isOk
}

// IsError returns true if the Result contains an error
func (r Result[T]) IsError() bool {
	return !r.isOk
}

// === Value Access ===

// Unwrap returns the contained value, panicking if there's an error
func (r Result[T]) Unwrap() T {
	if !r.isOk {
		panic(fmt.Sprintf("called Unwrap on error Result: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the contained value or a default value if there's an error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// UnwrapOrElse returns the contained value or computes a default using a function
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.isOk {
		return r.value
	}
	return fn(r.err)
}

// Expect returns the contained value, panicking with a custom message if there's an error
func (r Result[T]) Expect(msg string) T {
	if !r.isOk {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}
	return r.value
}

// === Error Access ===

// Err returns the contained error or nil if there's no error
func (r Result[T]) Err() error {
	if r.isOk {
		return nil
	}
	return r.err
}

// === Functional Operations ===

// Map transforms the contained value using the provided function
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if !r.isOk {
		return r
	}
	return Ok(fn(r.value))
}

// MapTo transforms the contained value to a different type
func MapTo[T, U any](r Result[T], fn func(T) U) Result[U] {
	if !r.isOk {
		return Fail[U](r.err)
	}
	return Ok(fn(r.value))
}

// FlatMap chains Result-returning operations
func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	if !r.isOk {
		return r
	}
	return fn(r.value)
}

// FlatMapTo chains Result-returning operations with type transformation
func FlatMapTo[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if !r.isOk {
		return Fail[U](r.err)
	}
	return fn(r.value)
}

// Filter keeps the value only if the predicate returns true
func (r Result[T]) Filter(predicate func(T) bool, errorMsg string) Result[T] {
	if !r.isOk {
		return r
	}
	if predicate(r.value) {
		return r
	}
	return Fail[T](errors.New(errorMsg))
}

// === Side Effects ===

// Inspect allows you to inspect the contained value without consuming the Result
func (r Result[T]) Inspect(fn func(T)) Result[T] {
	if r.isOk {
		fn(r.value)
	}
	return r
}

// InspectErr allows you to inspect the contained error without consuming the Result
func (r Result[T]) InspectErr(fn func(error)) Result[T] {
	if !r.isOk {
		fn(r.err)
	}
	return r
}

func (r Result[T]) WithErr(err error) Result[T] {
	return Result[T]{err: err}
}

// === Error Handling ===

// MapErr transforms the error using the provided function
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.isOk {
		return r
	}
	return Fail[T](fn(r.err))
}

// OrElse returns this Result if it's Ok, otherwise returns the result of the function
func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return fn(r.err)
}

// === Combinators ===

// And returns other if this Result is Ok, otherwise returns this Result
func (r Result[T]) And(other Result[T]) Result[T] {
	if r.isOk {
		return other
	}
	return r
}

// Or returns this Result if it's Ok, otherwise returns other
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return other
}

// === Conversion ===

// Note: ToOption removed as Option[T] is not idiomatic in Go
// Use IsOk() and Unwrap() directly, or handle with Result[T] methods

// String provides a string representation of the Result
func (r Result[T]) String() string {
	if r.isOk {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Error(%v)", r.err)
}

// === Integration with existing error handling ===

// Must returns the value or panics with the error
func (r Result[T]) Must() T {
	return r.Unwrap()
}

// === Early Return Pattern ===

// CatchIntoResult provides a clean early return pattern for error handling
// Usage: if someResult.CatchIntoResult(&targetResult) { return }
// Returns true if an error was caught, false if the operation was successful
func CatchIntoResult[T, U any](source Result[T], target *Result[U]) bool {
	if source.isOk {
		return false
	}
	*target = Fail[U](source.err)
	return true
}

// CatchInto provides error handling with automatic Result creation
// Usage: if someResult.CatchInto(&targetResult) { return targetResult }
// This is a generic version that works with any Result type
func CatchInto[T, U any](source Result[T], target *Result[U]) bool {
	if source.isOk {
		return false
	}
	*target = Fail[U](source.err)
	return true
}

// === Catchable Interface Implementation ===

// Catch implements the Catchable interface for standard error handling
func (r Result[T]) Catch(errPtr *error) bool {
	if r.isOk {
		return false
	}
	if errPtr != nil {
		*errPtr = r.err
	}
	return true
}

// CatchErr implements the Catchable interface for Error type handling
func (r Result[T]) CatchErr(errPtr *Error) bool {
	if r.isOk {
		return false
	}
	if errPtr != nil {
		*errPtr = NewError(r.err)
	}
	return true
}
