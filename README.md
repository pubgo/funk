# Funk

Funk is a comprehensive Go utility library that provides enhanced error handling, result types, feature flags, and various helper functions to simplify Go development.

## Features

### 🚀 Enhanced Error Handling
- Rich error wrapping with stack traces
- Error codes compatible with gRPC status codes
- Detailed error context and metadata
- Panic recovery mechanisms

### 📦 Result Types
- Functional programming-inspired Result and Error types
- Safe error handling without explicit nil checks
- Comprehensive API for mapping, filtering, and transforming results
- Async support with Future types

### 🎛️ Feature Flags
- Dynamic configuration at runtime
- Environment variable integration
- CLI flag generation
- Type-safe feature access

### 📝 Configuration Management
- YAML-based configuration with environment variable substitution
- Support for configuration merging and extension
- Expression engine for dynamic configuration values (similar to GitHub Actions workflow syntax)
- Hot-reloading capabilities

### 🪵 Structured Logging
- High-performance logging based on zerolog
- Context-aware logging with automatic field injection
- Error detail capture with stack traces
- Modular logger support with namespacing

### 🌍 Environment Management
- Normalized environment variable access
- Type-safe environment variable retrieval
- .env file loading and parsing
- Environment variable expansion

### 📚 Stack Trace Analysis
- Runtime stack trace capture
- Caller identification and metadata extraction
- Performance-optimized stack operations
- Deep Go runtime integration

### 🔧 Utility Functions
- Generic helper functions for slices, maps, and comparisons
- Assertion utilities
- Path and file utilities
- String and formatting helpers

## Installation

```bash
go get github.com/pubgo/funk/v2
```

## Quick Start

### Error Handling
```go
import "github.com/pubgo/funk/v2/errors"

err := errors.New("something went wrong", errors.Tags{"component": "database"})
err = errors.Wrap(err, "failed to connect")
errors.DebugPrint(err)
```

### Result Types
```go
import "github.com/pubgo/funk/v2/result"

// Create a successful result
res := result.OK(42)

// Create a failed result
errRes := result.Fail[int](errors.New("calculation failed"))

// Safely unwrap values
if value, ok := res.TryUnwrap(); ok {
    fmt.Printf("Got value: %d\n", value)
}

// Chain operations
result := result.OK(10).
    Map(func(x int) int { return x * 2 }).
    FlatMap(func(x int) result.Result[int] {
        if x > 0 {
            return result.OK(x + 1)
        }
        return result.Fail[int](errors.New("negative value"))
    })
```

### Feature Flags
```go
import "github.com/pubgo/funk/v2/features"

// Define a feature flag
var debugMode = features.Bool("debug", false, "Enable debug mode")

// Use the feature flag
if debugMode.Value() {
    log.Println("Debug mode enabled")
}
```

### Configuration
```go
import "github.com/pubgo/funk/v2/config"

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

cfg := config.Load[Config]()
fmt.Printf("Server will run on %s:%d\n", cfg.T.Server.Host, cfg.T.Server.Port)
```

### Logging
```go
import "github.com/pubgo/funk/v2/log"

logger := log.GetLogger("myapp")
logger.Info().Msg("Application started")
logger.Err(someError).Msg("An error occurred")
```

### Environment Variables
```go
import "github.com/pubgo/funk/v2/env"

// Get environment variable with fallback
host := env.GetOr("SERVER_HOST", "localhost")

// Type-safe environment variable access
port := env.GetInt("SERVER_PORT", "PORT")
debug := env.GetBool("DEBUG")
```

### Stack Trace Analysis
```go
import "github.com/pubgo/funk/v2/stack"

// Capture current stack trace
frames := stack.Trace()

// Get caller information
caller := stack.Caller(0)
fmt.Printf("Called from: %s:%d\n", caller.File, caller.Line)
```

## Modules

- **errors**: Enhanced error handling with wrapping, stack traces, and metadata
- **result**: Functional Result and Error types for safer error handling
- **features**: Feature flag system for runtime configuration
- **assert**: Assertion utilities for testing and validation
- **config**: Configuration management with multiple sources
- **log**: Enhanced logging capabilities
- **env**: Environment variable management
- **stack**: Stack trace analysis and caller identification

## Documentation

For detailed documentation, please visit:
- [Error Handling](./errors/README.md)
- [Result Types](./result/README.md)
- [Feature Flags](./features/README.md)
- [Configuration](./config/README.md)
- [Logging](./log/README.md)
- [Environment](./env/README.md)
- [Stack](./stack/README.md)

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## License

MIT