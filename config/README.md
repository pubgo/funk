# Configuration Module

The configuration module provides a flexible system for managing application configuration with support for YAML files, environment variable substitution, and configuration merging.

## Features

- **YAML-based Configuration**: Declarative configuration using YAML files
- **Environment Variable Substitution**: Automatic replacement of placeholders with environment values using `${ENV:"default_value"}` syntax
- **Expression Engine**: Dynamic configuration values using expressions similar to GitHub Actions workflow syntax (`$env.ENV`)
- **Configuration Merging**: Support for base and override configurations
- **Hot Reloading**: Dynamic configuration updates without application restart

## Installation

```bash
go get github.com/pubgo/funk/v2/config
```

## Quick Start

### Basic Configuration

```go
import "github.com/pubgo/funk/v2/config"

type ServerConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

type AppConfig struct {
    Server ServerConfig `yaml:"server"`
    Debug  bool         `yaml:"debug"`
}

// Load configuration
cfg := config.Load[AppConfig]()

// Access configuration values
fmt.Printf("Server: %s:%d\n", cfg.T.Server.Host, cfg.T.Server.Port)
fmt.Printf("Debug mode: %t\n", cfg.T.Debug)
```

### Environment Variable Substitution

Configuration files support environment variable substitution using the `${ENV:"default_value"}` syntax:

```yaml
# config.yaml
server:
  host: ${SERVER_HOST:"localhost"}
  port: ${SERVER_PORT:8080}
database:
  url: postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}
```

### Expression Engine

Advanced configuration values can be computed using expressions similar to GitHub Actions workflow syntax:

```yaml
# config.yaml
app:
  version: ${{env.VERSION}}
  build_time: ${{env.BUILD_TIME}}
  config_dir: ${{config_dir()}}
  cert_data: ${{embed("cert.pem")}}
```

Available expression functions:
- `env.ENV_VAR`: Access environment variables (similar to GitHub Actions syntax)
- `config_dir()`: Get the configuration directory
- `embed(filename)`: Embed file content as base64

### Configuration Structure

```go
type Resources struct {
    Resources       []string `yaml:"resources"`
    PatchResources  []string `yaml:"patch_resources"`
    PatchEnvs       []string `yaml:"patch_envs"`
}
```

## Core Concepts

### Configuration Loading

The module provides a generic `Load[T]()` function that loads and parses configuration files:

```go
cfg := config.Load[AppConfig]()
```

This function:
1. Finds the configuration file in predefined locations
2. Loads and parses the YAML content
3. Processes environment variable substitutions using both `${ENV}` and `${env.ENV}` syntax
4. Merges with additional configuration files if specified
5. Returns a typed configuration structure

### Environment Variable Integration

Environment variables can be used in configuration files with optional default values:

```yaml
# Syntax: ${VAR_NAME:"default_value"}
setting1: ${SETTING_1:"default_value"}
setting2: ${REQUIRED_SETTING}  # No default, will error if not set
```

Expression syntax (GitHub Actions style):
```yaml
# Syntax: ${env.VAR_NAME}
setting1: ${env.SETTING_1}
setting2: ${env.REQUIRED_SETTING}
```

### Configuration Merging

The module supports merging multiple configuration files:

```yaml
# config.yaml
resources:
  - configs/database.yaml
  - configs/cache.yaml
patch_resources:
  - configs/local-overrides.yaml
```

### Expression Evaluation

Advanced configuration values can be computed using expressions:

```yaml
# config.yaml
app:
  version: ${env.VERSION}
  build_time: ${env.BUILD_TIME}
  config_dir: ${config_dir()}
  cert_data: ${embed("cert.pem")}
```

## Advanced Usage

### Custom Configuration Paths

```go
// Set custom configuration path
config.SetConfigPath("/path/to/custom/config.yaml")

// Load configuration
cfg := config.Load[AppConfig]()
```

### Manual Configuration Loading

```go
cfg, err := config.LoadFromPath[AppConfig]("path/to/config.yaml")
if err != nil {
  panic(err)
}

// typed config value
_ = cfg.T

// export fully merged YAML for inspection/debugging
merged, err := config.LoadMergedConfigData("path/to/config.yaml")
if err != nil {
  panic(err)
}
fmt.Println(string(merged))
```

### Configuration Validation

```go
type ValidatedConfig struct {
    Port int `yaml:"port" validate:"min=1,max=65535"`
}

cfg := config.Load[ValidatedConfig]()
// Add custom validation logic here
```

## API Reference

### Core Functions

| Function                                | Description                                                        |
| --------------------------------------- | ------------------------------------------------------------------ |
| `Load[T]()`                             | Load and parse configuration generically                           |
| `TryLoad[T]()`                          | Try load typed configuration (returns error)                       |
| `LoadFromPath[T](cfgPath string)`       | Load typed configuration from specific path                        |
| `LoadMergedConfigData(cfgPath string)`  | Load fully merged processed config content (YAML bytes)            |
| `TryLoadMergedData()`                   | Try load merged processed config content via global path/discovery |
| `LoadMergedData()`                      | Load merged processed config content (panic on error)              |
| `ValidateEnvReferences(cfgPath string)` | Manually validate env references against `patch_envs`              |
| `SetConfigPath(path string)`            | Set custom configuration file path                                 |
| `GetConfigPath()`                       | Get current configuration file path                                |

### Configuration Structures

| Structure   | Description                                          |
| ----------- | ---------------------------------------------------- |
| `Cfg[T]`    | Wrapper containing loaded configuration and metadata |
| `Resources` | Configuration for resource loading and merging       |
| `Node`      | YAML node wrapper for flexible access                |

### Environment Integration

| Function                              | Description                          |
| ------------------------------------- | ------------------------------------ |
| `LoadEnvMap(cfgPath string)`          | Load environment configuration map   |
| `RegisterExpr(name string, expr any)` | Register custom expression functions |

## Best Practices

1. **Use Struct Tags**: Clearly define YAML mappings with struct tags
2. **Provide Defaults**: Specify sensible default values for configuration options
3. **Validate Early**: Implement configuration validation to catch errors early
4. **Environment Overrides**: Use environment variables for deployment-specific settings
5. **Modular Configs**: Split large configurations into logical modules
6. **Documentation**: Document configuration options and their purposes

## Integration Patterns

### With Feature Flags

```go
import (
    "github.com/pubgo/funk/v2/config"
    "github.com/pubgo/funk/v2/features"
)

type AppConfig struct {
    Debug bool `yaml:"debug"`
}

cfg := config.Load[AppConfig]()
debugMode := features.Bool("debug", cfg.T.Debug, "Enable debug mode")
```

### With Logging

```go
import (
    "github.com/pubgo/funk/v2/config"
    "github.com/pubgo/funk/v2/log"
)

type LogConfig struct {
    Level string `yaml:"level"`
}

cfg := config.Load[LogConfig]()
logger := log.GetLogger("app")
logger.Info().Str("config_level", cfg.T.Level).Msg("Configuration loaded")
```