# Errors Module

The errors module provides enhanced error handling capabilities for Go applications, offering rich context, stack traces, and metadata while maintaining compatibility with the standard `error` interface.

## Features

- **Rich Error Context**: Attach metadata and tags to errors
- **Stack Traces**: Automatic stack trace capture for debugging
- **Error Wrapping**: Preserve error chains with additional context
- **gRPC Compatibility**: Integration with gRPC status codes
- **Panic Recovery**: Safe recovery from panics with context preservation

## Installation

```bash
go get github.com/pubgo/funk/v2/errors
```

## Quick Start

### Creating Errors

```go
import "github.com/pubgo/funk/v2/errors"

// Simple error creation
err := errors.New("database connection failed")

// Error with metadata
err := errors.New("database connection failed", errors.Tags{
    "component": "database",
    "host": "localhost:5432",
})
```

### Error Wrapping

```go
// Wrap with additional context
err = errors.Wrap(err, "failed to initialize service")

// Wrap with formatted message
err = errors.Wrapf(err, "failed to initialize service after %d attempts", retryCount)

// Wrap with stack trace
err = errors.WrapStack(err)
```

### Error Inspection

```go
// Standard error checking
if err != nil {
    // Handle error
}

// Type assertion
if dbErr, ok := errors.AsA[*DatabaseError](err); ok {
    // Handle database-specific error
}

// Error chain traversal
for err != nil {
    fmt.Println(err.Error())
    err = errors.Unwrap(err)
}
```

## Core Concepts

### Error Types

The module provides two primary error types:

1. **`Err`**: Base error implementation with message, detail, and tags
2. **`ErrWrap`**: Wrapper that adds context and optionally stack traces

### Error Creation Functions

- `New(msg string, tags ...Tags)`: Create a new error with optional metadata
- `Errorf(format string, args ...any)`: Create a formatted error
- `Wrap(err error, msg string)`: Wrap an error with additional context
- `Wrapf(err error, format string, args ...any)`: Wrap with formatted context
- `WrapStack(err error)`: Wrap with stack trace information

### Error Inspection

- `As(err error, target any)`: Type assertion with error chain traversal
- `Is(err, target error)`: Check if error chain contains a specific error
- `Unwrap(err error)`: Retrieve the underlying error

## Advanced Usage

### Custom Error Types

```go
type DatabaseError struct {
    Op string
    Table string
    Err error
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database error in %s.%s: %v", e.Op, e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}
```

### Error Metadata

```go
// Attach metadata to errors
err := errors.New("payment processing failed", errors.Tags{
    "transaction_id": "txn_12345",
    "amount": 99.99,
    "currency": "USD",
})

// Access metadata
if tags := errors.GetTags(err); tags != nil {
    txnID := tags["transaction_id"]
    // Process transaction ID
}
```

### Stack Trace Analysis

```go
// Capture stack trace
err := errors.WrapStack(errors.New("critical failure"))

// Print detailed error information
errors.DebugPrint(err)

// Extract structured error data
data := errors.ErrJsonify(err)
```

## Integration

### With gRPC

```go
import "github.com/pubgo/funk/v2/errors/errcode"

// Convert to gRPC status
status := errcode.ConvertErr2Status(err)

// Create error code
code := &errorpb.ErrCode{
    Code: 500,
    StatusCode: errorpb.Code_Internal,
    Name: "INTERNAL_ERROR",
    Message: "Internal server error",
}
```

### With Logging

```go
import "github.com/pubgo/funk/v2/log"

// Log errors with context
log.Error().Err(err).Msg("operation failed")

// Structured error logging
log.Error().RawJSON("error", errors.JsonPrint(err)).Send()
```

## Best Practices

1. **Always Wrap Errors**: Provide context when passing errors up the call stack
2. **Use Stack Traces Sparingly**: Only in critical paths or for debugging
3. **Attach Relevant Metadata**: Include information that aids debugging
4. **Preserve Error Chains**: Use `Unwrap()` to maintain error causality
5. **Handle Errors at Appropriate Levels**: Don't catch and rethrow without adding value

## API Reference

### Error Creation

| Function | Description |
|----------|-------------|
| `New(msg string, tags ...Tags)` | Create new error with optional tags |
| `Errorf(format string, args ...any)` | Create formatted error |
| `Wrap(err error, msg string)` | Wrap error with context |
| `Wrapf(err error, format string, args ...any)` | Wrap with formatted context |
| `WrapStack(err error)` | Wrap with stack trace |
| `WrapTags(err error, tags Tags)` | Wrap with metadata tags |

### Error Inspection

| Function | Description |
|----------|-------------|
| `As(err error, target any)` | Type assertion with chain traversal |
| `Is(err, target error)` | Check error chain for specific error |
| `Unwrap(err error)` | Get underlying error |
| `GetErrorId(err error)` | Get unique error identifier |

### Utility Functions

| Function | Description |
|----------|-------------|
| `DebugPrint(err error)` | Pretty print error with stack trace |
| `JsonPrint(err error)` | Serialize error to JSON |
| `ErrJsonify(err error)` | Convert error to structured data |