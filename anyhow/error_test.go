package anyhow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pubgo/funk/anyhow"
)

func TestErrorBasicUsage(t *testing.T) {
	// Test successful Error (no error)
	okError := anyhow.ErrorOf(nil)
	if !okError.IsOk() {
		t.Fatal("Expected Ok Error")
	}
	if okError.IsError() {
		t.Fatal("Expected Ok Error")
	}

	// Test error Error
	err := fmt.Errorf("test error")
	errorError := anyhow.ErrorOf(err)
	if !errorError.IsError() {
		t.Fatal("Expected Error with error")
	}
	if errorError.IsOk() {
		t.Fatal("Expected Error with error")
	}

	if errorError.Err() == nil {
		t.Fatal("Expected underlying error")
	}
}

func TestErrorFunctionalOperations(t *testing.T) {
	// Test Map on error
	originalErr := fmt.Errorf("original error")
	errorError := anyhow.ErrorOf(originalErr)

	mappedError := errorError.Map(func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})

	if !strings.Contains(mappedError.Err().Error(), "wrapped:") {
		t.Fatalf("Expected wrapped error, got %v", mappedError.Err())
	}

	// Test Map on Ok Error
	okError := anyhow.ErrorOf(nil)
	mappedOk := okError.Map(func(err error) error {
		return fmt.Errorf("should not be called")
	})

	if !mappedOk.IsOk() {
		t.Fatal("Map on Ok Error should remain Ok")
	}
}

func TestErrorFlatMap(t *testing.T) {
	// Test FlatMap on error
	originalErr := fmt.Errorf("original error")
	errorError := anyhow.ErrorOf(originalErr)

	flatMapped := errorError.FlatMap(func(err error) anyhow.Error {
		return anyhow.ErrorFromString("new error from flatmap")
	})

	if flatMapped.Err().Error() != "new error from flatmap" {
		t.Fatalf("Expected new error, got %v", flatMapped.Err())
	}

	// Test FlatMap on Ok Error
	okError := anyhow.ErrorOf(nil)
	flatMappedOk := okError.FlatMap(func(err error) anyhow.Error {
		return anyhow.ErrorFromString("should not be called")
	})

	if !flatMappedOk.IsOk() {
		t.Fatal("FlatMap on Ok Error should remain Ok")
	}
}

func TestErrorOrElse(t *testing.T) {
	// Test OrElse on error
	originalErr := fmt.Errorf("original error")
	errorError := anyhow.ErrorOf(originalErr)

	recovered := errorError.OrElse(func(err error) anyhow.Error {
		return anyhow.ErrorOf(nil) // Recover to Ok
	})

	if !recovered.IsOk() {
		t.Fatal("OrElse should have recovered the error")
	}

	// Test OrElse on Ok Error
	okError := anyhow.ErrorOf(nil)
	orElseOk := okError.OrElse(func(err error) anyhow.Error {
		return anyhow.ErrorFromString("should not be called")
	})

	if !orElseOk.IsOk() {
		t.Fatal("OrElse on Ok Error should remain Ok")
	}
}

func TestErrorInspect(t *testing.T) {
	var inspectedError error

	// Test Inspect on error
	originalErr := fmt.Errorf("test error")
	errorError := anyhow.ErrorOf(originalErr)

	result := errorError.Inspect(func(err error) {
		inspectedError = err
	})

	if inspectedError == nil || inspectedError.Error() != "test error" {
		t.Fatalf("Expected to inspect error, got %v", inspectedError)
	}

	// Should return the same Error for chaining
	if !result.IsError() {
		t.Fatal("Inspect should return the same Error")
	}

	// Test Inspect on Ok Error
	inspectedError = nil
	okError := anyhow.ErrorOf(nil)
	okError.Inspect(func(err error) {
		inspectedError = err
	})

	if inspectedError != nil {
		t.Fatal("Inspect on Ok Error should not call the function")
	}
}

func TestErrorConversion(t *testing.T) {
	// Test ToResult conversion
	okError := anyhow.ErrorOf(nil)
	result := anyhow.ToResult(okError, "success value")

	if !result.IsOk() || result.Unwrap() != "success value" {
		t.Fatal("ToResult should convert Ok Error to Ok Result")
	}

	// Test error ToResult conversion
	errorError := anyhow.ErrorOf(fmt.Errorf("test error"))
	errorResult := anyhow.ToResult(errorError, "should not be used")

	if !errorResult.IsError() {
		t.Fatal("ToResult should convert Error with error to Error Result")
	}

	// Note: ToOption removed as Option[T] is not idiomatic in Go
	// Use ToResult instead for error handling
}

func TestErrorBackwardCompatibility(t *testing.T) {
	// Test WithErr (old API)
	originalErr := fmt.Errorf("original")
	errorError := anyhow.ErrorOf(originalErr)

	withErrResult := errorError.WithErr(func(err error) error {
		return fmt.Errorf("processed: %w", err)
	})

	if !strings.Contains(withErrResult.Err().Error(), "processed:") {
		t.Fatalf("WithErr should process the error, got %v", withErrResult.Err())
	}

	// Test OnErr (old API)
	var onErrCalled bool
	errorError.OnErr(func(err error) {
		onErrCalled = true
	})

	if !onErrCalled {
		t.Fatal("OnErr should be called for error")
	}

	// Test IsErr (old API)
	if !errorError.IsErr() {
		t.Fatal("IsErr should return true for error")
	}

	// Test GetErr (old API)
	getErr := errorError.GetErr()
	if getErr == nil {
		t.Fatal("GetErr should return error")
	}
}

func TestErrorCatchInto(t *testing.T) {
	var capturedErr error

	// Test CatchInto with error
	errorError := anyhow.ErrorOf(fmt.Errorf("test error"))
	caught := errorError.CatchInto(&capturedErr)

	if !caught {
		t.Fatal("CatchInto should return true when there's an error")
	}
	if capturedErr == nil || capturedErr.Error() != "test error" {
		t.Fatalf("CatchInto should capture the error, got %v", capturedErr)
	}

	// Test CatchInto with Ok Error
	capturedErr = nil
	okError := anyhow.ErrorOf(nil)
	caught = okError.CatchInto(&capturedErr)

	if caught {
		t.Fatal("CatchInto should return false when there's no error")
	}
	if capturedErr != nil {
		t.Fatal("CatchInto should not set error when there's no error")
	}
}

func TestErrorStringRepresentation(t *testing.T) {
	// Test String method
	okError := anyhow.ErrorOf(nil)
	if okError.String() != "Ok" {
		t.Fatalf("Expected 'Ok', got %s", okError.String())
	}

	errorError := anyhow.ErrorOf(fmt.Errorf("test error"))
	// Just check that it starts with "Error(" since the wrapped error format may vary
	if !strings.HasPrefix(errorError.String(), "Error(") {
		t.Fatalf("Expected string to start with 'Error(', got %s", errorError.String())
	}

	// Test Error method (error interface)
	if okError.Error() != "" {
		t.Fatalf("Expected empty string for Ok Error, got %s", okError.Error())
	}

	// The error should contain "test error" somewhere in the message
	if !strings.Contains(errorError.Error(), "test error") {
		t.Fatalf("Expected error to contain 'test error', got %s", errorError.Error())
	}
}

func TestErrorConstructors(t *testing.T) {
	// Test NewError
	err := fmt.Errorf("test")
	newError := anyhow.NewError(err)
	if newError.Err() == nil {
		t.Fatal("NewError should create Error with error")
	}

	// Test ErrorFromString
	stringError := anyhow.ErrorFromString("string error")
	if stringError.Err().Error() != "string error" {
		t.Fatalf("Expected 'string error', got %s", stringError.Err().Error())
	}

	// Test ErrorFromFormat
	formatError := anyhow.ErrorFromFormat("formatted %s %d", "error", 42)
	expected := "formatted error 42"
	if formatError.Err().Error() != expected {
		t.Fatalf("Expected '%s', got %s", expected, formatError.Err().Error())
	}
}

func TestErrorIntegrationWithAnyhow(t *testing.T) {
	// Test ToError function
	err := fmt.Errorf("test error")
	anyhowError := anyhow.ToError(err)
	if !anyhowError.IsError() {
		t.Fatal("ToError should create Error with error")
	}

	// Test ErrorResult function
	value, errorResult := anyhow.ErrorResult(func() (string, error) {
		return "success", nil
	})
	if value != "success" || !errorResult.IsOk() {
		t.Fatal("ErrorResult should handle successful function")
	}

	_, errorResult = anyhow.ErrorResult(func() (string, error) {
		return "", fmt.Errorf("function error")
	})
	if !errorResult.IsError() {
		t.Fatal("ErrorResult should handle error function")
	}
}

func TestErrorChaining(t *testing.T) {
	// Test method chaining
	result := anyhow.ErrorOf(fmt.Errorf("original error")).
		Map(func(err error) error {
			return fmt.Errorf("step1: %w", err)
		}).
		Map(func(err error) error {
			return fmt.Errorf("step2: %w", err)
		}).
		Inspect(func(err error) {
			// Should be called
		}).
		OrElse(func(err error) anyhow.Error {
			return anyhow.ErrorFromString("recovered")
		})

	if result.Err().Error() != "recovered" {
		t.Fatalf("Expected 'recovered', got %s", result.Err().Error())
	}
}
