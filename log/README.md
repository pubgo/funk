# Logging Module

The logging module provides a high-performance, structured logging system based on zerolog with support for context-aware logging, error detail capture, and modular logger management.

## Features

- **High-Performance Logging**: Based on zerolog for minimal overhead
- **Structured Logging**: JSON-formatted logs with automatic field extraction
- **Context-Aware Logging**: Request-scoped logging with automatic field injection
- **Error Detail Capture**: Automatic error enrichment with stack traces and metadata
- **Modular Loggers**: Namespaced loggers for different components
- **Multiple Output Formats**: Console and JSON output formats

## Design Notes

- The package intentionally uses `zerolog` as its concrete backend and exposes `Event` as an alias of `zerolog.Event` for a fluent, low-overhead API.
- `Logger.WithFields(...)` returns a derived logger and does not mutate the original logger.
- `WithLogger(ctx, ...)` is the convenience entry point for attaching a logger to a context.
- `WithFields(ctx, ...)` is the convenience entry point for attaching request-scoped fields to a context.
- `CreateFieldsCtx(...)` snapshots the initial field map so later mutations to the caller's map do not leak into logging context.
- `WithLogger(...)`, `WithFields(...)`, `GetFieldsFromCtx(...)`, and `WithDisabled(...)` are nil-safe helpers for request-scoped logging flows.
- Global helpers such as `log.Info(ctx)` and `log.Err(err, ctx)` honor the logger stored in the context.
- `log.Err(err, ctx)` automatically adds `error_id`, `error_chain`, `error_tags`, and `error_detail` for funk errors.
- `FromCtx(ctx)` is the ergonomic way to fetch the effective logger from a context when you need to keep chaining logger methods.
- `UpdateFieldsCtx(...)` returns a new context and does not mutate field maps already stored in parent contexts.
- When logger fields and context fields use the same key, **context fields win**. This lets request-scoped metadata override module defaults safely.

## Recommended Patterns

- **Typical request flow**: use `WithLogger` + `WithFields`, then call `log.Info(ctx)` / `log.Err(err, ctx)`.
- **Need direct logger chaining from context**: use `log.FromCtx(ctx)`.
- **Low-level explicit construction**: keep `CreateCtx` / `CreateFieldsCtx` for places where you want stricter, non-helper semantics.

## Installation

```bash
go get github.com/pubgo/funk/v2/log
```

## Quick Start

### Basic Logging

```go
import "github.com/pubgo/funk/v2/log"

// Use the global logger
log.Info().Msg("Application started")

// Log with fields
log.Info().Str("component", "server").Int("port", 8080).Msg("Server listening")

// Error logging with automatic error detail capture
err := someOperation()
if err != nil {
    log.Err(err).Msg("Operation failed")
}
```

### Modular Logging

```go
import "github.com/pubgo/funk/v2/log"

// Create a module-specific logger
logger := log.GetLogger("database")

// Use the module logger
logger.Info().Msg("Connecting to database")
logger.Err(dbError).Msg("Database connection failed")
```

### Context-Aware Logging

```go
import (
    "context"
    "github.com/pubgo/funk/v2/log"
)

// Create a context with fields
ctx := context.Background()
ctx = log.WithLogger(ctx, log.GetLogger("request"))
ctx = log.WithFields(ctx, map[string]any{
    "request_id": "abc123",
    "user_id": 42,
})

// Log with context (fields are automatically included)
log.Info(ctx).Msg("Processing request")

// Global helpers also use the context-bound logger automatically.

// If you need logger chaining, fetch the effective logger from context.
log.FromCtx(ctx).WithName("db").Info(ctx).Msg("Query finished")
```

### Context Helper Cheat Sheet

- `WithLogger`: bind a request-scoped logger to context
- `WithFields`: append request-scoped fields
- `log.Info(ctx)` / `log.Err(err, ctx)`: the common path for emitting logs
- `FromCtx`: fetch the effective logger when you need more chaining
- `CreateCtx` / `CreateFieldsCtx`: lower-level explicit constructors

## Core Concepts

### Logger Interface

The module provides a `Logger` interface for consistent logging operations:

```go
type Logger interface {
    Debug(ctx ...context.Context) *zerolog.Event
    Info(ctx ...context.Context) *zerolog.Event
    Warn(ctx ...context.Context) *zerolog.Event
    Error(ctx ...context.Context) *zerolog.Event
    Err(err error, ctx ...context.Context) *zerolog.Event
    // ... other methods
}
```

### Global Logger

A global logger is available for simple use cases:

```go
// Direct logging functions
log.Debug().Msg("Debug message")
log.Info().Msg("Info message")
log.Warn().Msg("Warning message")
log.Error().Msg("Error message")
log.Fatal().Msg("Fatal message")
log.Panic().Msg("Panic message")
```

### Module Loggers

Module-specific loggers provide namespacing and component isolation:

```go
// Create a logger for a specific module
dbLogger := log.GetLogger("database")
cacheLogger := log.GetLogger("cache")

// Use module-specific loggers
dbLogger.Info().Msg("Database operation")
cacheLogger.Info().Msg("Cache operation")
```

### Context Integration

Context-aware logging automatically includes relevant fields:

```go
// Add fields to context
ctx = log.WithFields(ctx, map[string]any{
    "trace_id": generateTraceID(),
    "span_id": generateSpanID(),
})

// Log with context (fields are automatically included)
log.Info(ctx).Msg("Processing with trace context")
```

## Advanced Usage

### Custom Logger Configuration

```go
import (
    "os"
    "github.com/rs/zerolog"
    "github.com/pubgo/funk/v2/log"
)

// Create a custom zerolog logger
customLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

// Set as global logger
log.SetLogger(&customLogger)
```

### Logger Levels

```go
import "github.com/rs/zerolog"

// Set global log level
zerolog.SetGlobalLevel(zerolog.InfoLevel)

// Create logger with specific level
logger := log.GetLogger("module").WithLevel(zerolog.DebugLevel)
```

### Error Detail Capture

The module automatically enriches error logs with detailed information:

```go
import "github.com/pubgo/funk/v2/errors"

// Create an error with metadata
err := errors.New("database connection failed", errors.Tags{
    "host": "localhost",
    "port": 5432,
})

// Log error (automatically includes error details)
log.Err(err).Msg("Connection attempt failed")
```

### Slog / Std Adapter Usage

```go
import (
    "log/slog"
    "github.com/pubgo/funk/v2/log"
)

ctx := log.WithLogger(nil, log.GetLogger("request"))

// slog adapter
slogger := slog.New(log.NewSlog(log.FromCtx(ctx)))
slogger.Info("handled request")

// std-style adapter
std := log.NewStd(log.FromCtx(ctx))
std.Println("request", "finished")
```

Notes:
- `NewSlog(nil)` and `NewStd(nil)` safely fall back to the global logger.
- `StdLogger.Println(...)` follows standard-library-style spacing between arguments.

### Event Building

```go
// Build complex log events
log.Info().
    Str("component", "auth").
    Int("user_id", 12345).
    Dict("metadata", log.NewEvent().
        Str("ip", "192.168.1.1").
        Str("user_agent", "Mozilla/5.0")).
    Msg("User login")
```

## API Reference

### Core Functions

| Function                                              | Description                                    |
| ----------------------------------------------------- | ---------------------------------------------- |
| `FromCtx(ctx context.Context)`                        | Get the effective logger from context          |
| `GetFromCtx(ctx context.Context, loggers ...Logger)`  | Low-level logger lookup with optional fallback |
| `GetLogger(names ...string)`                          | Get a module-specific logger                   |
| `SetLogger(log *zerolog.Logger)`                      | Set the global logger                          |
| `CreateCtx(ctx context.Context, ll Logger)`           | Create a context with an explicit logger       |
| `WithLogger(ctx context.Context, ll Logger)`          | Attach a logger to context                     |
| `WithFields(ctx context.Context, fields Fields)`      | Add fields to context                          |
| `CreateFieldsCtx(ctx context.Context, evt Fields)`    | Create a context with initial fields           |
| `UpdateFieldsCtx(ctx context.Context, fields Fields)` | Add or override fields in a derived context    |
| `GetFieldsFromCtx(ctx context.Context)`               | Extract fields from context                    |

### Logging Functions

| Function                                 | Description                                |
| ---------------------------------------- | ------------------------------------------ |
| `Debug(ctx ...context.Context)`          | Start a debug level log event              |
| `Info(ctx ...context.Context)`           | Start an info level log event              |
| `Warn(ctx ...context.Context)`           | Start a warning level log event            |
| `Error(ctx ...context.Context)`          | Start an error level log event             |
| `Err(err error, ctx ...context.Context)` | Start an error log event with error detail |
| `Fatal(ctx ...context.Context)`          | Start a fatal log event                    |
| `Panic(ctx ...context.Context)`          | Start a panic log event                    |

### Utility Functions

| Function                                    | Description                                    |
| ------------------------------------------- | ---------------------------------------------- |
| `NewEvent()`                                | Create a new dictionary event                  |
| `RecordErr(logs ...Logger)`                 | Create an error recording function             |
| `Output(w io.Writer)`                       | Create logger with custom output               |
| `OutputWriter(w func([]byte) (int, error))` | Create logger with custom writer function      |
| `NewSlog(log Logger)`                       | Create a `slog.Handler` backed by `log.Logger` |
| `NewStd(log Logger)`                        | Create a std-style logger adapter              |

## Best Practices

1. **Use Module Loggers**: Create module-specific loggers for better organization
2. **Include Context**: Use context-aware logging for request tracing
3. **Structure Logs**: Use structured logging with relevant fields
4. **Appropriate Levels**: Use appropriate log levels for different scenarios
5. **Avoid Sensitive Data**: Never log passwords, tokens, or other sensitive information
6. **Error Detail**: Always log errors with sufficient context for debugging
7. **Treat Loggers as Immutable**: Chain `WithName` / `WithFields` to create derived loggers instead of expecting in-place mutation
8. **Prefer Context For Request Data**: Put request-specific values in `context.Context` rather than module logger defaults

## Performance Notes

- The package keeps `zerolog` as the execution engine to preserve its allocation profile and fluent builder API.
- Recent regression tests ensure context-field merging does not leak state into reusable loggers.
- Benchmarks are included in `z_bench_test.go` for context merging and context-aware info logging so future changes can be measured quickly.

## Integration Patterns

### With Error Handling

```go
import (
    "github.com/pubgo/funk/v2/log"
    "github.com/pubgo/funk/v2/errors"
    "github.com/pubgo/funk/v2/result"
)

func processRequest() result.Error {
    err := someOperation()
    if err != nil {
        // Log error with full detail
        log.Err(err).Msg("Request processing failed")
        return result.ErrOf(err)
    }
    return result.Error{}
}
```

### With Configuration

```go
import (
    "github.com/pubgo/funk/v2/log"
    "github.com/pubgo/funk/v2/config"
)

type LogConfig struct {
    Level string `yaml:"level"`
    Format string `yaml:"format"`
}

cfg := config.Load[LogConfig]()
logger := log.GetLogger("app")

logger.Info().
    Str("log_level", cfg.T.Level).
    Str("log_format", cfg.T.Format).
    Msg("Logging configured")
```

### With Context and Tracing

```go
import (
    "context"
    "github.com/pubgo/funk/v2/log"
    "github.com/google/uuid"
)

func handleRequest(ctx context.Context) {
    // Add request ID to context
    requestID := uuid.New().String()
    ctx = log.WithLogger(ctx, log.GetLogger("request"))
    ctx = log.WithFields(ctx, map[string]any{
        "request_id": requestID,
    })
    
    // Log with request context
    log.Info(ctx).Msg("Handling request")

    // Continue chaining on the effective context logger when needed.
    log.FromCtx(ctx).WithName("downstream").Info(ctx).Msg("Calling processData")
    
    // Pass context to downstream functions
    processData(ctx)
}
```