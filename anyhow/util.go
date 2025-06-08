package anyhow

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

var errFnIsNil = errors.New("[fn] is nil")

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errFnIsNil)
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}

func try1[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errFnIsNil)
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return
}

func errMust(err error, args ...interface{}) {
	if generic.IsNil(err) {
		return
	}

	if len(args) > 0 {
		err = errors.Wrap(err, fmt.Sprint(args...))
	}

	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}

// TryContext executes a function with context cancellation support
func TryContext[T any](ctx context.Context, fn func(context.Context) (T, error)) Result[T] {
	if fn == nil {
		return Fail[T](errors.New("function is nil"))
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return Fail[T](ctx.Err())
	default:
	}

	resultChan := make(chan Result[T], 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					resultChan <- Fail[T](errors.Wrap(err, "panic in TryContext"))
				} else {
					resultChan <- Fail[T](errors.Errorf("panic in TryContext: %v", r))
				}
			}
		}()

		value, err := fn(ctx)
		resultChan <- From(value, err)
	}()

	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return Fail[T](ctx.Err())
	}
}

// BatchTry executes multiple functions concurrently and returns all results
func BatchTry[T any](fns ...func() (T, error)) []Result[T] {
	if len(fns) == 0 {
		return []Result[T]{}
	}

	results := make([]Result[T], len(fns))
	done := make(chan struct{})

	for i, fn := range fns {
		go func(index int, f func() (T, error)) {
			results[index] = Try(f)
			done <- struct{}{}
		}(i, fn)
	}

	// Wait for all goroutines to complete
	for range fns {
		<-done
	}

	return results
}

// RetryWith retries a function with exponential backoff
func RetryWith[T any](fn func() (T, error), attempts int) Result[T] {
	if attempts <= 0 {
		return Fail[T](errors.New("retry attempts must be positive"))
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		result := Try(fn)
		if result.IsOk() {
			return result
		}
		lastErr = result.Err()

		// Simple backoff - in production you might want exponential backoff
		if i < attempts-1 {
			// time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}
	}

	return Fail[T](errors.Wrapf(lastErr, "failed after %d attempts", attempts))
}

// Memoize caches the result of a function call
func Memoize[T any](fn func() (T, error)) func() Result[T] {
	var cached Result[T]
	var computed bool

	return func() Result[T] {
		if !computed {
			cached = Try(fn)
			computed = true
		}
		return cached
	}
}

// Collect transforms a slice using a Result-returning function
func Collect[T, U any](slice []T, fn func(T) Result[U]) Result[[]U] {
	results := make([]U, 0, len(slice))

	for _, item := range slice {
		result := fn(item)
		if result.IsError() {
			return Fail[[]U](result.Err())
		}
		results = append(results, result.Unwrap())
	}

	return Ok(results)
}

// CollectOptions transforms a slice into Results, filtering out errors
func CollectOptions[T, U any](slice []T, fn func(T) Result[U]) []U {
	var results []U

	for _, item := range slice {
		if res := fn(item); res.IsOk() {
			results = append(results, res.Unwrap())
		}
	}

	return results
}

// Partition separates a slice into two slices based on a predicate
func Partition[T any](slice []T, predicate func(T) bool) (truthy []T, falsy []T) {
	for _, item := range slice {
		if predicate(item) {
			truthy = append(truthy, item)
		} else {
			falsy = append(falsy, item)
		}
	}
	return
}

// FindFirst finds the first element matching a predicate
func FindFirst[T any](slice []T, predicate func(T) bool) Result[T] {
	for _, item := range slice {
		if predicate(item) {
			return Ok(item)
		}
	}
	return Fail[T](errors.New("no element found matching predicate"))
}

// === Error Collection ===

// MultiError represents multiple errors
type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("%d errors: ", len(m.Errors)))
	for i, err := range m.Errors {
		if i > 0 {
			msg.WriteString("; ")
		}
		msg.WriteString(err.Error())
	}
	return msg.String()
}

// CollectErrors collects all errors from a slice of Results
func CollectErrors[T any](results []Result[T]) Result[MultiError] {
	var errors []error

	for _, result := range results {
		if result.IsError() {
			errors = append(errors, result.Err())
		}
	}

	if len(errors) == 0 {
		return Fail[MultiError](fmt.Errorf("no errors found"))
	}

	return Ok(MultiError{Errors: errors})
}
