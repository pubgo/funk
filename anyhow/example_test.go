package anyhow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pubgo/funk/anyhow"
)

// Example data structures
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

func TestBasicResultUsage(t *testing.T) {
	// Test successful case
	result := anyhow.Ok(42)
	if !result.IsOk() {
		t.Fatal("Expected Ok result")
	}

	value := result.Unwrap()
	if value != 42 {
		t.Fatalf("Expected 42, got %d", value)
	}

	// Test error case
	errorResult := anyhow.Fail[int](fmt.Errorf("test error"))
	if !errorResult.IsError() {
		t.Fatal("Expected Error result")
	}

	err := errorResult.Err()
	if err == nil || err.Error() != "test error" {
		t.Fatalf("Expected test error, got %v", err)
	}
}

func TestFunctionalOperations(t *testing.T) {
	// Test Map
	result := anyhow.Ok(5).Map(func(x int) int { return x * 2 })
	if result.Unwrap() != 10 {
		t.Fatalf("Expected 10, got %d", result.Unwrap())
	}

	// Test MapTo (type transformation)
	stringResult := anyhow.MapTo(anyhow.Ok(42), func(x int) string {
		return fmt.Sprintf("Number: %d", x)
	})
	expected := "Number: 42"
	if stringResult.Unwrap() != expected {
		t.Fatalf("Expected %s, got %s", expected, stringResult.Unwrap())
	}

	// Test FlatMap
	chainResult := anyhow.Ok(5).FlatMap(func(x int) anyhow.Result[int] {
		if x > 0 {
			return anyhow.Ok(x * x)
		}
		return anyhow.Fail[int](fmt.Errorf("negative number"))
	})
	if chainResult.Unwrap() != 25 {
		t.Fatalf("Expected 25, got %d", chainResult.Unwrap())
	}
}

func TestErrorHandling(t *testing.T) {
	// Test MapErr
	errorResult := anyhow.Fail[string](fmt.Errorf("original error"))
	mappedError := errorResult.MapErr(func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})

	if !strings.Contains(mappedError.Err().Error(), "wrapped:") {
		t.Fatalf("Expected wrapped error, got %s", mappedError.Err().Error())
	}

	// Test OrElse
	recovered := errorResult.OrElse(func(err error) anyhow.Result[string] {
		return anyhow.Ok("recovered value")
	})
	if recovered.Unwrap() != "recovered value" {
		t.Fatalf("Expected recovered value, got %s", recovered.Unwrap())
	}
}

func TestFromFunction(t *testing.T) {
	// Test From with success
	result := anyhow.From("hello", nil)
	if !result.IsOk() || result.Unwrap() != "hello" {
		t.Fatal("From should handle nil error correctly")
	}

	// Test From with error
	result = anyhow.From("", fmt.Errorf("test error"))
	if !result.IsError() {
		t.Fatal("From should handle error correctly")
	}
}

func TestTryFunction(t *testing.T) {
	// Test Try with successful function
	result := anyhow.Try(func() (string, error) {
		return "success", nil
	})

	if !result.IsOk() || result.Unwrap() != "success" {
		t.Fatal("Try should handle successful function")
	}

	// Test Try with error function
	result = anyhow.Try(func() (string, error) {
		return "", fmt.Errorf("function error")
	})

	if !result.IsError() {
		t.Fatal("Try should handle error function")
	}
}

func TestInspectMethods(t *testing.T) {
	var inspectedValue int
	var inspectedError error

	// Test Inspect on Ok result
	anyhow.Ok(42).
		Inspect(func(v int) { inspectedValue = v }).
		InspectErr(func(err error) { inspectedError = err })

	if inspectedValue != 42 {
		t.Fatalf("Expected inspected value 42, got %d", inspectedValue)
	}
	if inspectedError != nil {
		t.Fatalf("Expected no error inspection, got %v", inspectedError)
	}

	// Reset and test Inspect on Error result
	inspectedValue = 0
	inspectedError = nil

	anyhow.Fail[int](fmt.Errorf("test error")).
		Inspect(func(v int) { inspectedValue = v }).
		InspectErr(func(err error) { inspectedError = err })

	if inspectedValue != 0 {
		t.Fatalf("Expected no value inspection, got %d", inspectedValue)
	}
	if inspectedError == nil || inspectedError.Error() != "test error" {
		t.Fatalf("Expected error inspection, got %v", inspectedError)
	}
}

func TestStandardLibraryIntegration(t *testing.T) {
	// Test JSON operations
	user := User{ID: 1, Name: "Alice", Age: 30}

	// Marshal to JSON
	jsonResult := anyhow.JsonMarshal(user)
	if jsonResult.IsError() {
		t.Fatalf("JSON marshal failed: %v", jsonResult.Err())
	}

	// Unmarshal back
	unmarshalResult := anyhow.JsonUnmarshal[User](jsonResult.Unwrap())
	if unmarshalResult.IsError() {
		t.Fatalf("JSON unmarshal failed: %v", unmarshalResult.Err())
	}

	unmarshaled := unmarshalResult.Unwrap()
	if unmarshaled.Name != "Alice" {
		t.Fatalf("Expected Alice, got %s", unmarshaled.Name)
	}

	// Test string parsing
	intResult := anyhow.Atoi("123")
	if intResult.IsError() || intResult.Unwrap() != 123 {
		t.Fatalf("Expected 123, got error or wrong value")
	}

	// Test error case
	errorIntResult := anyhow.Atoi("not-a-number")
	if !errorIntResult.IsError() {
		t.Fatal("Expected error for invalid number")
	}
}

// Note: TestOptionType removed as Option[T] is not idiomatic in Go
// Use Result[T] for error handling or standard Go patterns for optional values

func TestCombinators(t *testing.T) {
	// Test All - all succeed
	results := anyhow.All(
		anyhow.Ok(1),
		anyhow.Ok(2),
		anyhow.Ok(3),
	)

	if results.IsError() {
		t.Fatalf("All should succeed when all inputs succeed")
	}

	values := results.Unwrap()
	expected := []int{1, 2, 3}
	for i, v := range values {
		if v != expected[i] {
			t.Fatalf("Expected %v, got %v", expected, values)
		}
	}

	// Test All - one fails
	resultsWithError := anyhow.All(
		anyhow.Ok(1),
		anyhow.Fail[int](fmt.Errorf("test error")),
		anyhow.Ok(3),
	)

	if !resultsWithError.IsError() {
		t.Fatal("All should fail when any input fails")
	}
}

func TestRealWorldExample(t *testing.T) {
	// Simulate a real-world data processing pipeline
	jsonData := `{"users":[{"id":1,"name":"Alice","age":30},{"id":2,"name":"Bob","age":25}],"total":2}`

	transformedResult := anyhow.JsonUnmarshalFromString[UserResponse](jsonData).
		Map(func(response UserResponse) UserResponse {
			// Transform ages (add 1 year)
			for i := range response.Users {
				response.Users[i].Age++
			}
			return response
		})

	result := anyhow.FlatMapTo(transformedResult, func(response UserResponse) anyhow.Result[string] {
		// Filter adults only and create summary
		adults := []User{}
		for _, user := range response.Users {
			if user.Age >= 26 {
				adults = append(adults, user)
			}
		}

		if len(adults) == 0 {
			return anyhow.Fail[string](fmt.Errorf("no adults found"))
		}

		summary := fmt.Sprintf("Found %d adults", len(adults))
		return anyhow.Ok(summary)
	})

	if result.IsError() {
		t.Fatalf("Pipeline failed: %v", result.Err())
	}

	expected := "Found 2 adults"
	if result.Unwrap() != expected {
		t.Fatalf("Expected %s, got %s", expected, result.Unwrap())
	}
}

func TestChainedOperations(t *testing.T) {
	// Complex chaining example
	stringResult := anyhow.Ok("  hello world  ").
		Map(strings.TrimSpace).
		Map(strings.ToUpper).
		Filter(func(s string) bool { return len(s) > 5 }, "string too short")

	lengthResult := anyhow.FlatMapTo(stringResult, func(s string) anyhow.Result[int] {
		return anyhow.Ok(len(s))
	})

	result := lengthResult.Map(func(length int) int { return length * 2 })

	if result.IsError() {
		t.Fatalf("Chain failed: %v", result.Err())
	}

	// "HELLO WORLD" has 11 characters, * 2 = 22
	expected := 22
	if result.Unwrap() != expected {
		t.Fatalf("Expected %d, got %d", expected, result.Unwrap())
	}
}

// Benchmark to test performance
func BenchmarkResultOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := anyhow.Ok(i).
			Map(func(x int) int { return x * 2 }).
			Map(func(x int) int { return x + 1 }).
			Map(func(x int) int { return x / 2 })

		_ = result.Unwrap()
	}
}

func BenchmarkTraditionalErrorHandling(b *testing.B) {
	processValue := func(x int) (int, error) {
		x = x * 2
		x = x + 1
		x = x / 2
		return x, nil
	}

	for i := 0; i < b.N; i++ {
		result, err := processValue(i)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}
