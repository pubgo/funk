# Environment Module

The environment module provides utilities for working with environment variables, including normalized access, type-safe retrieval, and .env file loading.

## Features

- **Normalized Environment Access**: Consistent naming conventions for environment variables
- **Type-Safe Retrieval**: Safe conversion to common types (bool, int, float)
- **.env File Support**: Automatic loading of .env files with godotenv
- **Environment Expansion**: Variable substitution with envsubst
- **Flexible Lookup**: Multiple fallback mechanisms for variable resolution

## Installation

```bash
go get github.com/pubgo/funk/v2/env
```

## Quick Start

### Basic Environment Access

```go
import "github.com/pubgo/funk/v2/env"

// Get environment variable with fallback
host := env.GetOr("SERVER_HOST", "localhost")

// Get environment variable or panic if not set
port := env.MustGet("SERVER_PORT")

// Get with multiple fallback names
dbHost := env.Get("DB_HOST", "DATABASE_HOST", "MYSQL_HOST")
```

### Type-Safe Access

```go
import "github.com/pubgo/funk/v2/env"

// Boolean values (supports: true, 1, on, yes)
debug := env.GetBool("DEBUG", "ENABLE_DEBUG")

// Integer values
port := env.GetInt("PORT", "SERVER_PORT")

// Float values
timeout := env.GetFloat("TIMEOUT", "REQUEST_TIMEOUT")
```

### .env File Loading

```go
import "github.com/pubgo/funk/v2/env"

// Load .env file
env.LoadFiles(".env", "config/.env.local")

// Access variables from loaded files
apiKey := env.Get("API_KEY")
```

## Core Concepts

### Environment Variable Normalization

The module normalizes environment variable names to ensure consistency:

```go
// These are all equivalent after normalization:
// server-host => SERVER_HOST
// server.host => SERVER_HOST
// server/host => SERVER_HOST
// SERVER-HOST => SERVER_HOST

host := env.Get("server-host")  // Matches SERVER_HOST
```

### Multiple Name Lookup

Functions support multiple names for flexible variable resolution:

```go
// Looks for DB_HOST, then DATABASE_HOST, then MYSQL_HOST
dbHost := env.Get("DB_HOST", "DATABASE_HOST", "MYSQL_HOST")
```

### Safe Type Conversion

Type-safe getters provide graceful handling of conversion errors:

```go
// GetInt returns -1 on conversion error
port := env.GetInt("PORT")
if port == -1 {
    port = 8080  // default
}

// GetBool handles various truthy values
debug := env.GetBool("DEBUG")  // true for "true", "1", "on", "yes"
```

## Advanced Usage

### Environment Variable Expansion

```go
import "github.com/pubgo/funk/v2/env"

// Expand environment variables in strings
template := "http://${{HOST}}:${{PORT}}"
expanded := env.Expand(template).UnwrapOr(template)
```

### Environment Manipulation

```go
import "github.com/pubgo/funk/v2/env"

// Set environment variables
env.Set("CUSTOM_VAR", "value").Log()

// Delete environment variables
env.Delete("UNUSED_VAR").Log()

// Check if variable exists
if val, exists := env.Lookup("OPTIONAL_VAR"); exists {
    fmt.Printf("Optional var: %s\n", val)
}
```

### Environment Mapping

```go
import "github.com/pubgo/funk/v2/env"

// Get all environment variables as a map
allVars := env.Map()

// Process all variables
for key, value := range allVars {
    fmt.Printf("%s=%s\n", key, value)
}
```

## API Reference

### Core Functions

| Function | Description |
|----------|-------------|
| `Get(names ...string)` | Get environment variable with fallback names |
| `MustGet(names ...string)` | Get environment variable or panic if not set |
| `GetOr(name, defaultVal string)` | Get with explicit default value |
| `Set(key, value string)` | Set environment variable |
| `MustSet(key, value string)` | Set environment variable or panic on error |
| `Delete(key string)` | Delete environment variable |
| `MustDelete(key string)` | Delete environment variable or panic on error |

### Type-Safe Getters

| Function | Description |
|----------|-------------|
| `GetBool(names ...string)` | Get boolean value |
| `GetInt(names ...string)` | Get integer value |
| `GetFloat(names ...string)` | Get float value |
| `Lookup(key string)` | Lookup variable existence |

### File Operations

| Function | Description |
|----------|-------------|
| `LoadFiles(files ...string)` | Load .env files |
| `Expand(value string)` | Expand environment variables in string |

### Utility Functions

| Function | Description |
|----------|-------------|
| `Map()` | Get all environment variables as map |
| `Key(key string)` | Normalize environment variable key |
| `Normalize(key string)` | Normalize key with validity check |

## Best Practices

1. **Use Descriptive Names**: Choose clear, descriptive environment variable names
2. **Provide Sensible Defaults**: Always consider reasonable default values
3. **Validate Early**: Validate environment variables at startup
4. **Document Variables**: Document all environment variables your application uses
5. **Group Related Variables**: Use consistent naming prefixes for related variables
6. **Avoid Sensitive Logging**: Be careful not to log sensitive environment variables

## Integration Patterns

### With Configuration

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/config"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

// Environment variables are automatically substituted in config files
cfg := config.Load[Config]()

// Override with environment variables if needed
if host := env.Get("SERVER_HOST"); host != "" {
    cfg.T.Server.Host = host
}
```

### With Feature Flags

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/features"
)

// Create feature flag with environment variable default
var debugMode = features.Bool("debug", env.GetBool("DEBUG"), "Enable debug mode")

// Create feature flag with environment variable lookup
var logLevel = features.String("log_level", env.GetOr("LOG_LEVEL", "info"), "Logging level")
```

### With Logging

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/log"
)

logger := log.GetLogger("env")

// Log environment variable loading
env.LoadFiles(".env").Log(logger.RecordErr())

// Log specific variable access
host := env.Get("SERVER_HOST")
logger.Info().Str("host", host).Msg("Server host configured")
```