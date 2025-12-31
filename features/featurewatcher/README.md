# Feature Watcher Module

The Feature Watcher module provides file monitoring functionality that watches a specified directory for file changes and automatically parses and applies the file contents to feature flags.

## Features

- 📁 **Directory Monitoring**: Monitors all files in a specified directory
- 🔄 **Auto Reload**: Automatically reloads and applies changes when files are modified
- 📝 **Simple Format**: Supports `name=value` format configuration files
- ⚡ **Debounce**: Prevents frequent triggers for better performance
- 🔒 **Security Control**: Automatically skips non-mutable feature flags
- 📊 **Detailed Logging**: Records loading process and error information

## File Format

Configuration files use a simple `name=value` format, one feature flag per line:

```
# This is a comment, starting with # or //
# Empty lines are ignored

app_status=ok
replicas=3
debug=true
log_level=info
max_retries=5
```

## Quick Start

### Basic Usage

```go
package main

import (
	"log"
	"os"

	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/features/featurewatcher"
)

func main() {
	// Register feature flags
	_ = features.String("app_status", "ok", "Application status")
	_ = features.Int("replicas", 1, "Number of replicas")
	_ = features.Bool("debug", false, "Enable debug mode")

	// Create watch directory
	watchDir := "./features"
	os.MkdirAll(watchDir, 0755)

	// Create watcher
	watcher, err := featurewatcher.NewWatcher(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	// Start monitoring
	if err := watcher.Start(); err != nil {
		log.Fatal(err)
	}
	defer watcher.Stop()

	log.Println("Monitoring started, edit files in", watchDir, "to update feature flags")
	
	// Keep running
	select {}
}
```

### Custom Configuration

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithDebounce(1*time.Second), // Set debounce to 1 second
)
```

### Using Custom Feature Instance

```go
customFeature := features.NewFeature()
_ = customFeature.AddFunc("custom_flag", "Custom flag", value, nil)

watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithFeature(customFeature), // Use custom Feature instance
)
```

## Configuration Options

### WithFeature(f *features.Feature)

Specifies the Feature instance to use. If not specified, the global `defaultFeature` is used by default.

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithFeature(customFeature),
)
```

### WithDebounce(d time.Duration)

Sets the debounce time to prevent multiple triggers when files change frequently. Default is 500 milliseconds.

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithDebounce(1*time.Second),
)
```

## How It Works

1. **Initial Load**: Automatically loads all files in the directory on startup
2. **File Monitoring**: Uses `fsnotify` to monitor file changes in the directory
3. **Event Handling**: Triggers reload when files are written or created
4. **Debounce**: Only processes the last change within the debounce period
5. **Parse and Apply**: Parses file content, finds corresponding feature flags and updates values

## File Parsing Rules

1. **Format**: Each line must be in `name=value` format
2. **Comments**: Lines starting with `#` or `//` are ignored
3. **Empty Lines**: Empty lines are ignored
4. **Whitespace**: Spaces before and after names and values are automatically trimmed
5. **Error Handling**: Lines with invalid format are logged but don't interrupt processing

## Security Features

### Automatically Skip Non-Mutable Feature Flags

If a feature flag is marked as `mutable: false`, it won't be updated even if it appears in the file:

```go
_ = features.String("version", "1.0.0", "Version",
	map[string]any{
		"mutable": false, // Not mutable
	})
```

### Automatically Skip Unregistered Feature Flags

If a file contains an unregistered feature flag name, it will be logged but won't cause an error.

## Logging

The module outputs detailed log information:

- **INFO**: Monitor start/stop, successful file loading
- **DEBUG**: Feature flag update details
- **WARN**: Format errors, feature flags not found, etc.
- **ERROR**: File read errors, value setting failures, etc.

## Notes

1. **Directory Must Exist**: The monitored directory must exist, otherwise an error will be returned
2. **File Format**: Files must use `name=value` format, other formats will be ignored
3. **Feature Flags Must Be Registered**: Only registered feature flags will be updated
4. **Concurrency Safe**: File loading and feature flag updates are thread-safe
5. **Debounce Time**: Setting a reasonable debounce time can prevent frequent triggers but also increases latency

## Example Scenarios

### Scenario 1: Dynamic Configuration in Development

In development environments, you can dynamically adjust feature flags by modifying configuration files without restarting the application.

### Scenario 2: Multi-Environment Configuration

Different environments can use different configuration files, achieving dynamic configuration switching through file monitoring.

### Scenario 3: Configuration Center Integration

Can be integrated with file synchronization features of configuration centers (such as etcd, Consul) to achieve automatic configuration updates.

## Error Handling

The module handles the following error cases:

- Directory does not exist
- File read failures
- Format errors (logged as warnings, continues processing other lines)
- Feature flag does not exist (logged as debug info, skipped)
- Value setting failures (logged as errors, continues processing other lines)

All errors are logged and won't interrupt the monitoring process.

