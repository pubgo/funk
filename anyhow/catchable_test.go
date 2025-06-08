package anyhow

import (
	"fmt"
	"testing"
)

// TestResultCatchableInterface tests that Result implements Catchable correctly
func TestResultCatchableInterface(t *testing.T) {
	// Test Catch method with successful Result
	successResult := Ok(42)
	var stdErr error

	caught := successResult.Catch(&stdErr)
	if caught {
		t.Error("Catch should return false for successful Result")
	}
	if stdErr != nil {
		t.Error("Error should not be set for successful Result")
	}

	// Test Catch method with failed Result
	failResult := Fail[int](fmt.Errorf("test error"))
	var capturedStdErr error

	caught = failResult.Catch(&capturedStdErr)
	if !caught {
		t.Error("Catch should return true for failed Result")
	}
	if capturedStdErr == nil {
		t.Fatal("Error should be captured for failed Result")
	}
	if capturedStdErr.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", capturedStdErr.Error())
	}

	// Test CatchErr method with successful Result
	var anyhowErr Error
	caught = successResult.CatchErr(&anyhowErr)
	if caught {
		t.Error("CatchErr should return false for successful Result")
	}
	if anyhowErr.IsError() {
		t.Error("Error should not be set for successful Result")
	}

	// Test CatchErr method with failed Result
	var capturedAnyhowErr Error
	caught = failResult.CatchErr(&capturedAnyhowErr)
	if !caught {
		t.Error("CatchErr should return true for failed Result")
	}
	if !capturedAnyhowErr.IsError() {
		t.Fatal("Error should be captured for failed Result")
	}
	if capturedAnyhowErr.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", capturedAnyhowErr.Error())
	}
}

// TestErrorCatchableInterface tests that Error implements Catchable correctly
func TestErrorCatchableInterface(t *testing.T) {
	// Test Catch method with successful Error (no error)
	okError := ErrorOf(nil)
	var stdErr error

	caught := okError.Catch(&stdErr)
	if caught {
		t.Error("Catch should return false for Ok Error")
	}
	if stdErr != nil {
		t.Error("Error should not be set for Ok Error")
	}

	// Test Catch method with failed Error
	failError := ErrorOf(fmt.Errorf("test error"))
	var capturedStdErr error

	caught = failError.Catch(&capturedStdErr)
	if !caught {
		t.Error("Catch should return true for failed Error")
	}
	if capturedStdErr == nil {
		t.Fatal("Error should be captured for failed Error")
	}
	if capturedStdErr.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", capturedStdErr.Error())
	}

	// Test CatchErr method with successful Error (no error)
	var anyhowErr Error
	caught = okError.CatchErr(&anyhowErr)
	if caught {
		t.Error("CatchErr should return false for Ok Error")
	}
	if anyhowErr.IsError() {
		t.Error("Error should not be set for Ok Error")
	}

	// Test CatchErr method with failed Error
	var capturedAnyhowErr Error
	caught = failError.CatchErr(&capturedAnyhowErr)
	if !caught {
		t.Error("CatchErr should return true for failed Error")
	}
	if !capturedAnyhowErr.IsError() {
		t.Fatal("Error should be captured for failed Error")
	}
	if capturedAnyhowErr.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", capturedAnyhowErr.Error())
	}
}

// TestCatchableInterfacePolymorphism tests polymorphic usage of Catchable interface
func TestCatchableInterfacePolymorphism(t *testing.T) {
	// Function that works with any Catchable
	handleCatchable := func(c Catchable) (hasStdErr bool, hasAnyhowErr bool) {
		var stdErr error
		var anyhowErr Error

		stdCaught := c.Catch(&stdErr)
		anyhowCaught := c.CatchErr(&anyhowErr)

		return stdCaught, anyhowCaught
	}

	// Test with successful Result
	successResult := Ok("success")
	stdCaught, anyhowCaught := handleCatchable(successResult)
	if stdCaught || anyhowCaught {
		t.Error("No errors should be caught for successful Result")
	}

	// Test with failed Result
	failResult := Fail[string](fmt.Errorf("result error"))
	stdCaught, anyhowCaught = handleCatchable(failResult)
	if !stdCaught || !anyhowCaught {
		t.Error("Both error types should be caught for failed Result")
	}

	// Test with Ok Error
	okError := ErrorOf(nil)
	stdCaught, anyhowCaught = handleCatchable(okError)
	if stdCaught || anyhowCaught {
		t.Error("No errors should be caught for Ok Error")
	}

	// Test with failed Error
	failError := ErrorOf(fmt.Errorf("error type error"))
	stdCaught, anyhowCaught = handleCatchable(failError)
	if !stdCaught || !anyhowCaught {
		t.Error("Both error types should be caught for failed Error")
	}
}

// TestCatchableWithNilPointers tests Catchable interface with nil pointers
func TestCatchableWithNilPointers(t *testing.T) {
	failResult := Fail[int](fmt.Errorf("test error"))
	failError := ErrorOf(fmt.Errorf("test error"))

	// Test Result with nil pointers (should not panic)
	caught := failResult.Catch(nil)
	if !caught {
		t.Error("Should return true even with nil pointer")
	}

	caught = failResult.CatchErr(nil)
	if !caught {
		t.Error("Should return true even with nil pointer")
	}

	// Test Error with nil pointers (should not panic)
	caught = failError.Catch(nil)
	if !caught {
		t.Error("Should return true even with nil pointer")
	}

	caught = failError.CatchErr(nil)
	if !caught {
		t.Error("Should return true even with nil pointer")
	}
}

// TestCatchableCompileTimeInterface ensures types implement the interface at compile time
func TestCatchableCompileTimeInterface(t *testing.T) {
	// These assignments ensure the types implement Catchable at compile time
	var _ Catchable = Ok(42)
	var _ Catchable = Fail[int](fmt.Errorf("error"))
	var _ Catchable = ErrorOf(fmt.Errorf("error"))
	var _ Catchable = ErrorOf(nil)

	t.Log("All types successfully implement Catchable interface")
}
