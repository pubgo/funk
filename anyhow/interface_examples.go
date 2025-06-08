package anyhow

import (
	"fmt"
	"log"
)

// === Generic Functions Using Interfaces ===

// CheckStatus is a generic function that works with any Checkable type
func CheckStatus(c Checkable) string {
	if c.IsOk() {
		return "✓ Success"
	}
	return "✗ Failed: " + c.String()
}

// SafeUnwrap safely unwraps any Unwrappable type with a default value
func SafeUnwrap[T any](u Unwrappable[T], defaultValue T) T {
	return u.UnwrapOr(defaultValue)
}

// LogError logs error information from any ErrorAccessible type
func LogError(ea ErrorAccessible, context string) {
	if err := ea.Err(); err != nil {
		log.Printf("[%s] Error occurred: %v", context, err)
	}
}

// === Polymorphic Operations ===

// ProcessCheckable demonstrates polymorphic behavior with different types
func ProcessCheckable(items ...Checkable) (successes, failures int) {
	for _, item := range items {
		if item.IsOk() {
			successes++
		} else {
			failures++
		}
	}
	return
}

// BatchUnwrap unwraps multiple Unwrappable values with defaults
func BatchUnwrap[T any](items []Unwrappable[T], defaultValue T) []T {
	results := make([]T, len(items))
	for i, item := range items {
		results[i] = item.UnwrapOr(defaultValue)
	}
	return results
}

// === Error Collection Utilities ===

// ErrorCollector collects errors from multiple ErrorAccessible types
type ErrorCollector struct {
	errors []error
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{errors: make([]error, 0)}
}

func (ec *ErrorCollector) Collect(ea ErrorAccessible) *ErrorCollector {
	if err := ea.Err(); err != nil {
		ec.errors = append(ec.errors, err)
	}
	return ec
}

func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

func (ec *ErrorCollector) Errors() []error {
	return ec.errors
}

func (ec *ErrorCollector) JoinedError() error {
	if len(ec.errors) == 0 {
		return nil
	}
	return JoinErrors(ec.errors...)
}

// === Adapter Functions ===

// WrapAsCheckable wraps a boolean and string as a Checkable interface
type boolCheckable struct {
	ok   bool
	desc string
}

func (bc boolCheckable) IsOk() bool     { return bc.ok }
func (bc boolCheckable) IsError() bool  { return !bc.ok }
func (bc boolCheckable) String() string { return bc.desc }

func WrapBoolAsCheckable(ok bool, description string) Checkable {
	return boolCheckable{ok: ok, desc: description}
}

// === Higher-Order Functions ===

// MapUnwrappable applies a transformation to unwrappable values
func MapUnwrappable[T, U any](items []Unwrappable[T], transform func(T) U, defaultValue T) []U {
	results := make([]U, len(items))
	for i, item := range items {
		value := item.UnwrapOr(defaultValue)
		results[i] = transform(value)
	}
	return results
}

// FilterCheckable filters Checkable items based on their status
func FilterCheckable(items []Checkable, keepOk bool) []Checkable {
	var filtered []Checkable
	for _, item := range items {
		if item.IsOk() == keepOk {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// === Validation Pipeline ===

// ValidationStep represents a single validation step
type ValidationStep[T any] interface {
	Validate(T) Checkable
}

// ValidationPipeline runs multiple validation steps
type ValidationPipeline[T any] struct {
	steps []ValidationStep[T]
}

func NewValidationPipeline[T any]() *ValidationPipeline[T] {
	return &ValidationPipeline[T]{steps: make([]ValidationStep[T], 0)}
}

func (vp *ValidationPipeline[T]) AddStep(step ValidationStep[T]) *ValidationPipeline[T] {
	vp.steps = append(vp.steps, step)
	return vp
}

func (vp *ValidationPipeline[T]) Validate(value T) []Checkable {
	results := make([]Checkable, len(vp.steps))
	for i, step := range vp.steps {
		results[i] = step.Validate(value)
	}
	return results
}

func (vp *ValidationPipeline[T]) IsValid(value T) bool {
	for _, step := range vp.steps {
		if step.Validate(value).IsError() {
			return false
		}
	}
	return true
}

// === Usage Examples ===

// ExampleUsage demonstrates how to use the interfaces in practice
func ExampleUsage() {
	// Create different types
	result := Ok(42)
	errorResult := ErrorOf(fmt.Errorf("something went wrong"))

	// Use generic functions
	fmt.Println(CheckStatus(result))      // ✓ Success
	fmt.Println(CheckStatus(errorResult)) // ✗ Failed: Error(...)

	// Safe unwrapping
	value := SafeUnwrap(result, 0) // 42

	fmt.Printf("Unwrapped value: %d\n", value)

	// Collect errors
	collector := NewErrorCollector()
	collector.Collect(result).Collect(errorResult)

	if collector.HasErrors() {
		fmt.Printf("Collected %d errors\n", len(collector.Errors()))
	}
}
