# Result Module

The result module provides functional programming-inspired error handling through `Result[T]` and `Error` types. It eliminates explicit nil checks and enables safe, composable error handling with method chaining.

## Features

- **Type-Safe Error Handling**: Eliminate nil-check bugs with generic Result types
- **Method Chaining**: Fluent API for transforming and combining operations
- **Pattern Matching**: Exhaustive error handling with Match functions
- **Async Support**: Future types for asynchronous computations
- **Standard Compatibility**: Works with existing Go error handling patterns

## Installation

```bash
go get github.com/pubgo/funk/v2/result
```

## Quick Start

### Basic Usage

```go
import "github.com/pubgo/funk/v2/result"

// Create successful result
success := result.OK(42)

// Create failed result
failure := result.Fail[int](errors.New("calculation failed"))

// Safe unwrapping
if value, ok := success.TryUnwrap(); ok {
    fmt.Printf("Value: %d\n", value)
}

// Method chaining
processed := result.OK("hello").
    Map(strings.ToUpper).
    FlatMap(func(s string) result.Result[string] {
        if len(s) > 0 {
            return result.OK(s + " world")
        }
        return result.Fail[string](errors.New("empty string"))
    })
```

### Error Handling

```go
// Pattern matching
result.OK(42).Match(
    func(value int) { fmt.Printf("Success: %d\n", value) },
    func(err error) { fmt.Printf("Error: %v\n", err) }
)

// Error-only operations
result.ErrOf(errors.New("something went wrong")).
    Log().
    InspectErr(func(err error) { 
        // Additional error processing
    })
```

## Core Concepts

### Result[T]

`Result[T]` represents either a successful value of type `T` or an error. It provides a safe alternative to returning `(T, error)` pairs.

#### Creation Functions

- `OK(v T)`: Create successful result
- `Fail[T](err error)`: Create failed result
- `Wrap(v T, err error)`: Create result from value/error pair
- `WrapFn(fn func() (T, error))`: Create result from function

#### Transformation Methods

- `Map(func(T) T)`: Transform successful value
- `FlatMap(func(T) Result[T])`: Transform with potential for new errors
- `Validate(func(T) error)`: Validate value with potential error
- `Or(default T)`: Provide fallback value for errors

#### Consumption Methods

- `Unwrap()`: Get value or panic on error
- `TryUnwrap()`: Safely get value with success indicator
- `UnwrapOr(default T)`: Get value or default on error
- `Expect(msg string)`: Get value or panic with custom message

### Error

`Error` is a specialized result type for error-only operations, providing utilities for error handling and logging.

#### Creation Functions

- `ErrOf(err error)`: Create from error
- `ErrOfFn(fn func() error)`: Create from error-returning function
- `Errorf(format string, args ...any)`: Create formatted error

#### Error Handling Methods

- `Log()`: Log error with stack trace
- `InspectErr(func(error))`: Process error without consuming it
- `Match(func(), func(error))`: Pattern match on error presence

### Future[T]

`Future[T]` represents asynchronous computations that will eventually produce a `Result[T]`.

#### Creation Functions

- `Async(fn func() Result[T])`: Create from async function
- `AsyncErr(fn func() Error)`: Create from async error function

#### Usage

```go
// Start async computation
future := result.Async(func() result.Result[int] {
    // Simulate work
    time.Sleep(time.Second)
    return result.OK(42)
})

// Wait for result
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result := future.Await(ctx)
```

## Advanced Usage

### Chaining Operations

```go
// Complex transformation pipeline
result := getUser(id).
    FlatMap(validateUser).
    Map(encryptUserData).
    FlatMap(saveUser).
    Inspect(func(savedUser User) {
        log.Info().Msgf("User saved: %s", savedUser.ID)
    }).
    MapErr(func(err error) error {
        return fmt.Errorf("failed to process user: %w", err)
    })
```

### Working with Collections

```go
// Process multiple results
results := []result.Result[int]{
    result.OK(1),
    result.OK(2),
    result.Fail[int](errors.New("conversion failed")),
    result.OK(4),
}

// Collect all successful values or first error
collected := result.Collect(results)
values, err := collected.UnwrapErr()

// Partition into successes and failures
successes, failures := result.Partition(results)
```

### Pattern Matching

```go
// Exhaustive error handling
result.OK(42).Match(
    func(value int) {
        fmt.Printf("Success: %d\n", value)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    }
)

// Result-producing pattern matching
result.OK(42).MatchWithResult(
    func(value int) result.Result[string] {
        return result.OK(strconv.Itoa(value))
    },
    func(err error) result.Result[string] {
        return result.Fail[string](err)
    }
)
```

## Integration Patterns

### With Standard Error Handling

```go
// Convert from standard (T, error) pattern
func divide(a, b float64) result.Result[float64] {
    if b == 0 {
        return result.Fail[float64](errors.New("division by zero"))
    }
    return result.OK(a / b)
}

// Convert to standard (T, error) pattern
value, err := divide(10, 2).UnwrapErr()
```

### With Context

```go
// Context-aware operations
func fetchUser(ctx context.Context, id string) result.Result[User] {
    user, err := database.GetUser(ctx, id)
    return result.Wrap(user, err)
}

// Async with context
future := result.Async(func() result.Result[User] {
    return fetchUser(ctx, userID)
})

result := future.Await(ctx)
```

## Best Practices

1. **Prefer Result[T] over (T, error)**: For better composability and fewer nil checks
2. **Use TryUnwrap for Safe Extraction**: Avoid panics in uncertain situations
3. **Chain Operations for Readability**: Use Map/FlatMap for data transformation pipelines
4. **Handle Errors Appropriately**: Use Match for exhaustive error handling
5. **Leverage Async for IO-bound Operations**: Use Future[T] for non-blocking computations

## API Reference

### Result[T] Creation

| Function | Description |
|----------|-------------|
| `OK(v T)` | Create successful result |
| `Fail[T](err error)` | Create failed result |
| `Wrap(v T, err error)` | Create from value/error pair |
| `WrapFn(fn func() (T, error))` | Create from function |

### Result[T] Inspection

| Method | Description |
|--------|-------------|
| `IsOK() bool` | Check if result is successful |
| `IsErr() bool` | Check if result is error |
| `TryUnwrap() (T, bool)` | Safe value extraction |

### Result[T] Transformation

| Method | Description |
|--------|-------------|
| `Map(func(T) T)` | Transform successful value |
| `FlatMap(func(T) Result[T])` | Transform with potential new errors |
| `Validate(func(T) error)` | Validate with potential error |

### Result[T] Consumption

| Method | Description |
|--------|-------------|
| `Unwrap() T` | Get value or panic |
| `UnwrapOr(default T) T` | Get value or default |
| `Expect(msg string) T` | Get value or panic with message |

### Error Creation

| Function | Description |
|----------|-------------|
| `ErrOf(err error)` | Create from error |
| `ErrOfFn(fn func() error)` | Create from error function |
| `Errorf(format string, args ...any)` | Create formatted error |

### Error Handling

| Method | Description |
|--------|-------------|
| `Log()` | Log error with context |
| `InspectErr(func(error))` | Process error without consuming |
| `Match(func(), func(error))` | Pattern match on error presence |