package anyhow

import (
	"context"

	"github.com/pubgo/funk/errors"
)

// === Integration with Standard Go Error Handling ===

// CatchPanic recovers from panics and converts them to Results
func CatchPanic[T any](fn func() T) Result[T] {
	defer func() {
		if r := recover(); r != nil {
			// This will be handled by the Result returned below
		}
	}()

	// We need to handle panics differently since we can't return from defer
	// Let's use a channel-based approach for safety
	resultChan := make(chan Result[T], 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					resultChan <- Fail[T](errors.Wrap(err, "panic recovered"))
				} else {
					resultChan <- Fail[T](errors.Errorf("panic recovered: %v", r))
				}
			}
		}()

		value := fn()
		resultChan <- Ok(value)
	}()

	return <-resultChan
}

// TryPanic is similar to Try but also catches panics
func TryPanic[T any](fn func() (T, error)) Result[T] {
	return CatchPanic(func() T {
		value, err := fn()
		if err != nil {
			panic(err)
		}
		return value
	})
}

// === Multiple Result Handling ===

// All returns Ok with all values if all Results are Ok, otherwise returns the first error
func All[T any](results ...Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, result := range results {
		if result.IsError() {
			return Fail[[]T](result.err)
		}
		values = append(values, result.value)
	}
	return Ok(values)
}

// Any returns the first Ok result, or the last error if all are errors
func Any[T any](results ...Result[T]) Result[T] {
	var lastErr error
	for _, result := range results {
		if result.IsOk() {
			return result
		}
		lastErr = result.err
	}
	return Fail[T](lastErr)
}

// === Utility Functions for Common Patterns ===

// IfThen creates a Result based on a condition
func IfThen[T any](condition bool, value T, err error) Result[T] {
	if condition {
		return Ok(value)
	}
	return Fail[T](err)
}

// Note: WhenSome removed as Option[T] is not idiomatic in Go
// Use Result[T] methods directly for error handling

// FirstOk returns the first Ok result, or an error if none found
func FirstOk[T any](results []Result[T]) Result[T] {
	for _, result := range results {
		if result.IsOk() {
			return result
		}
	}
	return Fail[T](errors.New("no successful results found"))
}

// === Context Integration ===

// WithContext adds context to error handling (for compatibility with existing patterns)
func WithContext[T any](ctx context.Context, result Result[T]) Result[T] {
	if result.IsError() {
		// Add context information to the error
		err := errors.WrapKV(result.err, "context", ctx)
		return Fail[T](err)
	}
	return result
}

// === Error Type Integration ===

// ToError converts a standard error to anyhow.Error
func ToError(err error) Error {
	return ErrorOf(err)
}

// ErrorResult creates an Error and a Result[T] pair from a function
func ErrorResult[T any](fn func() (T, error)) (T, Error) {
	value, err := fn()
	return value, ErrorOf(err)
}

// === Legacy Compatibility (Simplified) ===

// Must panics if the error is not nil (for quick migration from old code)
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Check returns true if error is not nil (for conditional error handling)
func Check(err error) bool {
	return err != nil
}

// === Deprecated: Legacy functions for migration compatibility ===
// These functions are kept for backward compatibility but should be avoided in new code

// RecoveryErr - DEPRECATED: Use Try or CatchPanic instead
func RecoveryErr(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				*setter = err
			} else {
				*setter = errors.Errorf("panic: %v", r)
			}

			for _, callback := range callbacks {
				*setter = callback(*setter)
				if *setter == nil {
					return
				}
			}
		}
	}()
}

// Wrap - DEPRECATED: Use From instead
func Wrap[T any](v T, err error) Result[T] {
	return From(v, err)
}

// WrapFn - DEPRECATED: Use Try instead
func WrapFn[T any](fn func() (T, error)) Result[T] {
	return Try(fn)
}
