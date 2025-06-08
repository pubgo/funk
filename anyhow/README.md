# Anyhow - Functional Error Handling for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/pubgo/funk/anyhow.svg)](https://pkg.go.dev/github.com/pubgo/funk/anyhow)
[![Go Report Card](https://goreportcard.com/badge/github.com/pubgo/funk/anyhow)](https://goreportcard.com/report/github.com/pubgo/funk/anyhow)

**Anyhow** is a library that brings Rust-inspired functional error handling to Go, reducing the need for repetitive `if err != nil` checks while maintaining type safety and performance.

## 🌟 Features

- **🔗 Functional Chaining**: Chain operations with `Map`, `FlatMap`, `Filter`
- **🛡️ Type Safety**: Full generic support with compile-time type checking
- **⚡ High Performance**: Zero-allocation design with direct value storage
- **📦 Standard Library Integration**: Built-in wrappers for common Go operations
- **🎯 Go-Idiomatic**: Follows Go conventions and patterns
- **🔄 Combinators**: Work with multiple Results elegantly

## 🚀 Quick Start

```go
import "github.com/pubgo/funk/anyhow"

// Traditional Go error handling
func traditionalWay() (string, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return "", err
    }
    
    var config Config
    err = json.Unmarshal(data, &config)
    if err != nil {
        return "", err
    }
    
    if config.Name == "" {
        return "", errors.New("name is required")
    }
    
    return strings.ToUpper(config.Name), nil
}

// With Anyhow - Functional style
func functionalWay() anyhow.Result[string] {
    return anyhow.ReadFile("config.json").
        FlatMapTo(anyhow.JsonUnmarshal[Config]).
        Filter(func(c Config) bool { return c.Name != "" }, "name is required").
        Map(func(c Config) string { return strings.ToUpper(c.Name) })
}
```

## 📋 API Overview

### Core Types

```go
// Result[T] - Contains either a value of type T or an error
type Result[T any] struct { /* ... */ }

// Error - Represents an operation that may have failed (for error-only operations)
type Error struct { /* ... */ }
```

> **Note**: We removed `Option[T]` as it's not idiomatic in Go. Go has zero values and pointers for optional values.

### Constructors

```go
// Create successful Result
result := anyhow.Ok(42)

// Create failed Result
result := anyhow.Fail[int](errors.New("something went wrong"))

// From Go's (value, error) pattern
result := anyhow.From(strconv.Atoi("123"))

// Safely execute function
result := anyhow.Try(func() (int, error) {
    return riskyOperation()
})

// Create Error type for error-only operations
errorResult := anyhow.ErrorOf(fmt.Errorf("something went wrong"))
```

### Functional Operations

#### Map - Transform Values
```go
result := anyhow.Ok(5).Map(func(x int) int { 
    return x * 2 
})  // Ok(10)
```

#### FlatMap - Chain Operations
```go
result := anyhow.Ok("123").FlatMap(func(s string) anyhow.Result[int] {
    return anyhow.Atoi(s)
})  // Ok(123)
```

#### Filter - Conditional Values
```go
result := anyhow.Ok(42).Filter(func(x int) bool { 
    return x > 0 
}, "must be positive")  // Ok(42)
```

### Error Handling

```go
result := riskyOperation().
    MapErr(func(err error) error {
        return fmt.Errorf("operation failed: %w", err)
    }).
    OrElse(func(err error) anyhow.Result[string] {
        return anyhow.Ok("default value")
    })
```

### Side Effects

```go
result := anyhow.Ok(42).
    Inspect(func(value int) {
        log.Printf("Got value: %d", value)
    }).
    InspectErr(func(err error) {
        log.Printf("Got error: %v", err)
    })
```

## 🔧 Standard Library Integration

### JSON Operations
```go
// Marshal to JSON
jsonResult := anyhow.JsonMarshal(user)

// Unmarshal from JSON
userResult := anyhow.JsonUnmarshal[User](jsonData)

// Chain JSON operations
result := anyhow.JsonMarshal(data).
    Map(string).  // Convert to string
    FlatMapTo(anyhow.JsonUnmarshal[OtherType])
```

### File Operations
```go
// Read file
content := anyhow.ReadTextFile("data.txt")

// Write file
result := anyhow.WriteTextFile("output.txt", "content", 0644)

// Copy file
copied := anyhow.CopyFile("source.txt", "dest.txt")
```

### String Parsing
```go
// Parse integer
number := anyhow.Atoi("123")

// Parse float
float := anyhow.ParseFloat("3.14", 64)

// Parse boolean
boolean := anyhow.ParseBool("true")
```

## 🎯 Go-Idiomatic Optional Values

Instead of Option[T], use Go's natural patterns:

```go
// Use pointers for optional values
var optionalUser *User
if condition {
    optionalUser = &User{Name: "Alice"}
}

// Use zero values
var optionalString string  // "" means no value
var optionalInt int        // 0 means no value

// Use Result[T] for operations that might fail
userResult := anyhow.GetUser(id)  // Result[User] instead of Option[User]
if userResult.IsOk() {
    user := userResult.Unwrap()
    // use user
}
```

## 🚫 Error Type

For operations that only care about success/failure (no return value):

```go
// Create Error instances
okError := anyhow.ErrorOf(nil)                    // Success
errorError := anyhow.ErrorOf(fmt.Errorf("failed")) // Failure

// Check status
if okError.IsOk() {
    // Handle success
}

// Functional operations
processedError := errorError.
    Map(func(err error) error {
        return fmt.Errorf("processed: %w", err)
    }).
    OrElse(func(err error) anyhow.Error {
        return anyhow.ErrorOf(nil) // Recover
    })

// Convert to other types
result := anyhow.ToResult(okError, "success value")  // Ok("success value")
option := anyhow.ToOption(errorError, 42)            // None[int]()

// Backward compatibility with old API
var capturedErr error
if errorError.CatchInto(&capturedErr) {
    // Handle captured error
}
```

## 🔄 Combinators

### Work with Multiple Results

```go
// All must succeed
results := anyhow.All(
    anyhow.Ok(1),
    anyhow.Ok(2),
    anyhow.Ok(3),
)  // Ok([1, 2, 3])

// First success wins
result := anyhow.Any(
    anyhow.Fail[int](error1),
    anyhow.Ok(42),
    anyhow.Fail[int](error2),
)  // Ok(42)
```

### Collection Operations

```go
numbers := []string{"1", "2", "3", "4"}

// Transform all or fail
results := anyhow.Collect(numbers, anyhow.Atoi)  // Ok([1, 2, 3, 4])

// Filter successful operations
validNumbers := anyhow.CollectOptions(numbers, func(s string) anyhow.Option[int] {
    return anyhow.Atoi(s).ToOption()
})  // [1, 2, 3, 4]
```

## 📊 Performance

Anyhow is designed for performance with zero-allocation operations:

```
BenchmarkResultOperations-8           1000000000    0.52 ns/op    0 B/op    0 allocs/op
BenchmarkTraditionalErrorHandling-8   1000000000    0.33 ns/op    0 B/op    0 allocs/op
```

## 🌍 Real-World Examples

### HTTP API Processing
```go
func ProcessUserRequest(userID string) anyhow.Result[UserResponse] {
    return anyhow.NotEmpty(userID).
        FlatMapTo(func(id string) anyhow.Result[User] {
            return fetchUser(id)
        }).
        FlatMapTo(func(user User) anyhow.Result[UserResponse] {
            return anyhow.All(
                fetchUserProfile(user.ID),
                fetchUserPreferences(user.ID),
                fetchUserStats(user.ID),
            ).Map(func(results []interface{}) UserResponse {
                return UserResponse{
                    User:        user,
                    Profile:     results[0].(Profile),
                    Preferences: results[1].(Preferences),
                    Stats:       results[2].(Stats),
                }
            })
        }).
        Inspect(func(response UserResponse) {
            log.Printf("Successfully processed user: %s", response.User.Name)
        })
}
```

### Configuration Loading
```go
func LoadConfig(path string) anyhow.Result[Config] {
    return anyhow.ReadTextFile(path).
        FlatMapTo(anyhow.JsonUnmarshalFromString[Config]).
        Map(func(config Config) Config {
            config.SetDefaults()
            return config
        }).
        Filter(func(config Config) bool {
            return config.Validate()
        }, "invalid configuration").
        Inspect(func(config Config) {
            log.Printf("Loaded config from %s", path)
        })
}
```

### Data Pipeline
```go
func ProcessData(input []string) anyhow.Result[ProcessedData] {
    return anyhow.Ok(input).
        Map(func(data []string) []string {
            return filter(data, isValid)
        }).
        Filter(func(data []string) bool {
            return len(data) > 0
        }, "no valid data found").
        FlatMapTo(func(data []string) anyhow.Result[ProcessedData] {
            return anyhow.Collect(data, processItem).
                Map(combineResults)
        })
}
```

## 📚 Migration Guide

Migrating from traditional Go error handling or the old Anyhow API? Check out our [Migration Guide](MIGRATION.md) for detailed instructions and examples.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by Rust's `Result<T, E>` and `Option<T>` types
- Thanks to the Go generics proposal for making this possible
- Built on top of the excellent error handling patterns in the Go community 