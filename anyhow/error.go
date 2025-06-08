package anyhow

import (
	"fmt"

	"github.com/pubgo/funk/errors"
)

// Error represents an operation that may have failed
// It implements multiple interfaces for functional programming patterns
type Error struct {
	err error
}

// Ensure Error implements all relevant interfaces at compile time
var _ Checkable = (*Error)(nil)
var _ ErrorAccessible = (*Error)(nil)
var _ Catchable = (*Error)(nil)

// === Constructors ===

// NewError creates a new Error from an error
func NewError(err error) Error {
	if err == nil {
		return Error{}
	}
	return Error{err: errors.WrapCaller(err, 1)}
}

// ErrorOf creates an Error from an error (alias for NewError)
func ErrorOf(err error) Error {
	return NewError(err)
}

// ErrorFromString creates an Error from a string message
func ErrorFromString(message string) Error {
	return Error{err: errors.New(message)}
}

// ErrorFromFormat creates an Error from a format string and arguments
func ErrorFromFormat(format string, args ...interface{}) Error {
	return Error{err: errors.Errorf(format, args...)}
}

// === State Checking ===

// IsOk returns true if there's no error
func (e Error) IsOk() bool {
	return e.err == nil
}

// IsError returns true if there's an error (alias for !IsOk for consistency)
func (e Error) IsError() bool {
	return e.err != nil
}

// === Error Access ===

// Err returns the underlying error
func (e Error) Err() error {
	return e.err
}

// GetErr returns the underlying error (backward compatibility)
func (e Error) GetErr() error {
	if e.err == nil {
		return nil
	}
	return errors.WrapCaller(e.err, 1)
}

// Unwrap returns the underlying error (for error unwrapping)
func (e Error) Unwrap() error {
	return e.err
}

// === Functional Operations ===

// Map transforms the error if present, otherwise returns the same Error
func (e Error) Map(fn func(error) error) Error {
	if e.err == nil {
		return e
	}
	return Error{err: fn(e.err)}
}

// FlatMap chains Error-returning operations
func (e Error) FlatMap(fn func(error) Error) Error {
	if e.err == nil {
		return e
	}
	return fn(e.err)
}

// OrElse returns this Error if it's Ok, otherwise returns the result of the function
func (e Error) OrElse(fn func(error) Error) Error {
	if e.err == nil {
		return e
	}
	return fn(e.err)
}

// === Side Effects ===

// Inspect allows you to inspect the error without consuming the Error
func (e Error) Inspect(fn func(error)) Error {
	if e.err != nil {
		fn(e.err)
	}
	return e
}

// === Conversion ===

// ToResult converts the Error to a Result[T]
func ToResult[T any](e Error, value T) Result[T] {
	if e.err == nil {
		return Ok(value)
	}
	return Fail[T](e.err)
}

// Note: ToOption removed as Option[T] is not idiomatic in Go
// Use ToResult instead for error handling

// === Panic Operations ===

// Must panics if there's an error
func (e Error) Must() {
	if e.err != nil {
		panic(e.err)
	}
}

// Expect panics with a custom message if there's an error
func (e Error) Expect(msg string) {
	if e.err != nil {
		panic(fmt.Sprintf("%s: %v", msg, e.err))
	}
}

// === String Representation ===

// String provides a string representation of the Error
func (e Error) String() string {
	if e.err == nil {
		return "Ok"
	}
	return fmt.Sprintf("Error(%v)", e.err)
}

// Error implements the error interface
func (e Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// === Backward Compatibility Functions ===

// WithErr transforms the error using the provided function (old API)
func (e Error) WithErr(fn func(error) error) Error {
	if e.err == nil {
		return e
	}
	return Error{err: fn(e.err)}
}

// OnErr allows you to inspect the error (old API)
func (e Error) OnErr(fn func(error)) {
	if e.err != nil {
		fn(e.err)
	}
}

// IsErr returns true if there's an error (old API, alias for IsError)
func (e Error) IsErr() bool {
	return e.err != nil
}

// CatchInto sets the error into the provided error pointer if this Error has an error
func (e Error) CatchInto(errPtr *error) bool {
	if e.err == nil {
		return false
	}
	if errPtr != nil {
		*errPtr = e.err
	}
	return true
}

// === Catchable Interface Implementation ===

// Catch implements the Catchable interface for standard error handling
func (e Error) Catch(errPtr *error) bool {
	if e.err == nil {
		return false
	}
	if errPtr != nil {
		*errPtr = e.err
	}
	return true
}

// CatchErr implements the Catchable interface for Error type handling
func (e Error) CatchErr(errPtr *Error) bool {
	if e.err == nil {
		return false
	}
	if errPtr != nil {
		*errPtr = e
	}
	return true
}

// === Utility Error Functions ===

// ErrorWithStack wraps an error with stack trace information
func ErrorWithStack(err error) error {
	if err == nil {
		return nil
	}
	return errors.WrapStack(err)
}

// ErrorWithMessage wraps an error with additional message
func ErrorWithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, message)
}

// JoinErrors combines multiple errors into one
func JoinErrors(errs ...error) error {
	var nonNilErrors []error
	for _, err := range errs {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}

	if len(nonNilErrors) == 0 {
		return nil
	}

	if len(nonNilErrors) == 1 {
		return nonNilErrors[0]
	}

	// Create a simple joined error
	var message string
	for i, err := range nonNilErrors {
		if i > 0 {
			message += "; "
		}
		message += err.Error()
	}

	return errors.New(fmt.Sprintf("multiple errors: %s", message))
}
