# Features Module

The features module provides a centralized system for managing application configuration through feature flags. It supports runtime configuration, environment variable integration, and CLI flag generation.

## Features

- **Centralized Configuration**: Single registry for all feature flags
- **Multiple Sources**: Environment variables, CLI flags, and programmatic configuration
- **Type Safety**: Strongly typed feature values (bool, string, int, float, json)
- **Metadata Support**: Tags for organizing and documenting features
- **CLI Integration**: Automatic generation of CLI flags for urfave/cli

## Installation

```bash
go get github.com/pubgo/funk/v2/features
```

## Quick Start

### Defining Feature Flags

```go
import "github.com/pubgo/funk/v2/features"

// Boolean feature flag
var debugMode = features.Bool("debug", false, "Enable debug mode")

// String feature flag with metadata
var logLevel = features.String("log_level", "info", "Logging level",
    map[string]any{
        "group": "logging",
        "allowed": []string{"debug", "info", "warn", "error"},
    })

// Integer feature flag
var maxRetries = features.Int("max_retries", 3, "Maximum retry attempts")

// JSON feature flag
type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}
var serverConfig = features.Json("server_config", &Config{}, "Server configuration")
```

### Using Feature Values

```go
// Access feature values
if debugMode.Value() {
    log.Println("Debug mode enabled")
}

// Update feature values programmatically
err := maxRetries.Set("5")
if err != nil {
    log.Printf("Failed to set max_retries: %v", err)
}

// Access with metadata
flag := features.Lookup("log_level")
if flag != nil {
    group := flag.Tags["group"]
    fmt.Printf("Log level feature is in group: %v\n", group)
}
```

## Core Concepts

### Feature Registry

The module uses a central registry pattern to manage all feature flags:

```go
// Default registry (most common usage)
var debugMode = features.Bool("debug", false, "Enable debug mode")

// Custom registry
customFeatures := features.NewFeature()
var customFlag = customFeatures.Bool("custom", false, "Custom feature")
```

### Value Types

The module supports several value types with appropriate parsing:

1. **BoolValue**: Boolean flags with flexible parsing ("true", "1", "on", "yes")
2. **StringValue**: String values
3. **IntValue**: Integer values
4. **FloatValue**: Floating-point values
5. **JsonValue**: JSON-marshallable values

### Metadata and Tags

Feature flags can include metadata through tags:

```go
var feature = features.Bool("feature", false, "Description",
    map[string]any{
        "group": "performance",
        "experimental": true,
        "owner": "team-name",
    })
```

## Advanced Usage

### Environment Variable Integration

Feature flags automatically integrate with environment variables:

```go
// Feature flag named "debug" will read from:
// - FEATURE_DEBUG (automatically generated)
// - DEBUG (direct mapping)

// Custom environment variable mapping
var custom = features.String("custom", "default", "Description")
// Reads from FEATURE_CUSTOM
```

### CLI Flag Generation

Generate CLI flags for urfave/cli applications:

```go
import "github.com/pubgo/funk/v2/features/featureflags"

app := &cli.App{
    Flags: append(
        baseFlags,
        featureflags.GetFlags()..., // Add all feature flags as CLI options
    ),
    Action: func(c *cli.Context) error {
        // Feature values are automatically updated from CLI
        if debugMode.Value() {
            log.Println("Debug mode enabled via CLI")
        }
        return nil
    },
}
```

### Feature Enumeration

Iterate over all registered features:

```go
features.VisitAll(func(flag *features.Flag) {
    fmt.Printf("Feature: %s = %v\n", flag.Name, flag.Value.Get())
    
    // Access metadata
    if group, ok := flag.Tags["group"]; ok {
        fmt.Printf("  Group: %v\n", group)
    }
})
```

### Dynamic Updates

Feature values can be updated at runtime:

```go
// Programmatic updates
err := debugMode.Set("true")
if err != nil {
    log.Printf("Failed to enable debug mode: %v", err)
}

// Updates from strings (useful for configuration loading)
values := map[string]string{
    "debug": "true",
    "max_retries": "5",
    "log_level": "debug",
}

for name, value := range values {
    flag := features.Lookup(name)
    if flag != nil {
        err := flag.Value.Set(value)
        if err != nil {
            log.Printf("Failed to set %s: %v", name, err)
        }
    }
}
```

## Integration Patterns

### With Configuration Files

Combine with configuration file loading:

```go
func loadConfigFromFile(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    
    var config map[string]any
    if err := json.Unmarshal(data, &config); err != nil {
        return err
    }
    
    // Update feature flags from config
    features.VisitAll(func(flag *features.Flag) {
        if value, exists := config[flag.Name]; exists {
            strValue := fmt.Sprintf("%v", value)
            if err := flag.Value.Set(strValue); err != nil {
                log.Printf("Failed to set %s from config: %v", flag.Name, err)
            }
        }
    })
    
    return nil
}
```

### With Observability

Integrate with monitoring systems:

```go
// Export feature values for monitoring
func exportFeatures() map[string]any {
    metrics := make(map[string]any)
    features.VisitAll(func(flag *features.Flag) {
        metrics[flag.Name] = flag.Value.Get()
    })
    return metrics
}

// Periodically export to metrics system
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        metrics := exportFeatures()
        // Send to metrics system (Prometheus, etc.)
        sendMetrics(metrics)
    }
}()
```

## Best Practices

1. **Define Features at Package Level**: For easy access throughout the application
2. **Use Descriptive Names and Usage Strings**: Make features self-documenting
3. **Apply Metadata for Organization**: Use tags to group related features
4. **Leverage Environment Variables**: Enable configuration without code changes
5. **Update Features Atomically**: When changing multiple related features
6. **Document Experimental Features**: Use metadata to mark unstable features
7. **Monitor Feature Usage**: Track which features are actually used

## API Reference

### Value Creation Functions

| Function | Description | Environment Variable |
|----------|-------------|---------------------|
| `Bool(name, default, usage, tags...)` | Boolean feature | FEATURE_NAME |
| `String(name, default, usage, tags...)` | String feature | FEATURE_NAME |
| `Int(name, default, usage, tags...)` | Integer feature | FEATURE_NAME |
| `Float(name, default, usage, tags...)` | Float feature | FEATURE_NAME |
| `Json[T](name, default, usage, tags...)` | JSON feature | FEATURE_NAME |

### Feature Registry

| Function | Description |
|----------|-------------|
| `Lookup(name string)` | Find feature by name |
| `VisitAll(func(*Flag))` | Iterate over all features |
| `NewFeature()` | Create custom registry |

### Value Methods

| Method | Description |
|--------|-------------|
| `Value() T` | Get current value |
| `Set(string) error` | Update from string |
| `String() string` | String representation |
| `Type() ValueType` | Get value type |
| `Get() any` | Get as interface{} |

### Flag Properties

| Property | Description |
|----------|-------------|
| `Name` | Feature name |
| `Usage` | Description string |
| `Value` | Current value wrapper |
| `Tags` | Metadata map |