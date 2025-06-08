package anyhow

import (
	"fmt"
	"testing"
)

func TestCheckableInterface(t *testing.T) {
	// Test with Result
	result := Ok(42)
	if !result.IsOk() {
		t.Fatal("Result should be Ok")
	}

	errorResult := Fail[int](fmt.Errorf("test error"))
	if !errorResult.IsError() {
		t.Fatal("Error result should be Error")
	}

	// Test with Error
	err := ErrorOf(fmt.Errorf("test"))
	if !err.IsError() {
		t.Fatal("Error should be Error")
	}

	okErr := ErrorOf(nil)
	if !okErr.IsOk() {
		t.Fatal("Ok Error should be Ok")
	}
}

func TestGenericFunctions(t *testing.T) {
	// Test CheckStatus with different types
	result := Ok(42)
	status := CheckStatus(result)
	if status != "✓ Success" {
		t.Fatalf("Expected success status, got %s", status)
	}

	err := ErrorOf(fmt.Errorf("test error"))
	status = CheckStatus(err)
	if status == "✓ Success" {
		t.Fatal("Error should not show success status")
	}

	// Test SafeUnwrap
	value := SafeUnwrap(result, 0)
	if value != 42 {
		t.Fatalf("Expected 42, got %d", value)
	}
}

func TestErrorCollector(t *testing.T) {
	collector := NewErrorCollector()

	// Collect from different types
	result := Fail[string](fmt.Errorf("result error"))
	err := ErrorOf(fmt.Errorf("error type error"))
	okResult := Ok("success")

	collector.Collect(result).Collect(err).Collect(okResult)

	if !collector.HasErrors() {
		t.Fatal("Collector should have errors")
	}

	errors := collector.Errors()
	if len(errors) != 2 {
		t.Fatalf("Expected 2 errors, got %d", len(errors))
	}

	joinedErr := collector.JoinedError()
	if joinedErr == nil {
		t.Fatal("Joined error should not be nil")
	}
}

func TestBoolCheckable(t *testing.T) {
	checkable := WrapBoolAsCheckable(true, "success")
	if !checkable.IsOk() {
		t.Fatal("Bool checkable should be Ok")
	}

	checkable = WrapBoolAsCheckable(false, "failed")
	if !checkable.IsError() {
		t.Fatal("Bool checkable should be Error")
	}
}

func TestBatchUnwrap(t *testing.T) {
	items := []Unwrappable[int]{
		Ok(1),
		Ok(2),
		Ok(3),
	}

	results := BatchUnwrap(items, 0)
	expected := []int{1, 2, 3}

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Fatalf("Expected %d at index %d, got %d", expected[i], i, result)
		}
	}
}

func TestProcessCheckable(t *testing.T) {
	items := []Checkable{
		Ok(42),
		Fail[string](fmt.Errorf("error")),
		ErrorOf(nil),
		ErrorOf(fmt.Errorf("another error")),
	}

	successes, failures := ProcessCheckable(items...)
	expectedSuccesses := 2
	expectedFailures := 2

	if successes != expectedSuccesses {
		t.Fatalf("Expected %d successes, got %d", expectedSuccesses, successes)
	}

	if failures != expectedFailures {
		t.Fatalf("Expected %d failures, got %d", expectedFailures, failures)
	}
}
