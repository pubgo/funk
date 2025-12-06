# Stack Module

The stack module provides utilities for stack trace analysis and caller identification, enabling detailed debugging and error reporting capabilities.

## Features

- **Runtime Stack Trace Capture**: Efficient capture of runtime stack traces
- **Caller Identification**: Accurate identification of calling functions
- **Stack Frame Analysis**: Detailed analysis of stack frames with metadata
- **Performance Optimization**: Cached stack trace processing for minimal overhead
- **Go Runtime Integration**: Deep integration with Go's runtime internals

## Installation

```bash
go get github.com/pubgo/funk/v2/stack
```

## Quick Start

### Basic Stack Trace Capture

```go
import "github.com/pubgo/funk/v2/stack"

// Capture current stack trace
frames := stack.Trace()

// Print stack frames
for _, frame := range frames {
    fmt.Printf("%s\n", frame.String())
}
```

### Caller Identification

```go
import "github.com/pubgo/funk/v2/stack"

// Get immediate caller
caller := stack.Caller(0)
fmt.Printf("Called from: %s:%d %s\n", caller.File, caller.Line, caller.Name)

// Get parent caller
parent := stack.Caller(1)
fmt.Printf("Parent caller: %s:%d %s\n", parent.File, parent.Line, parent.Name)
```

### Function-Based Caller Identification

```go
import "github.com/pubgo/funk/v2/stack"

func myFunction() {
    // Get caller of this function
    caller := stack.CallerWithFunc(myFunction)
    fmt.Printf("myFunction called from: %s\n", caller.String())
}
```

## Core Concepts

### Stack Frames

The module represents stack frames with detailed metadata:

```go
type Frame struct {
    Name string // Function name
    Pkg  string // Package path
    File string // File path
    Line int    // Line number
}
```

### Stack Trace Capture

Capture full stack traces for debugging:

```go
// Capture complete stack trace
frames := stack.Trace()

// Process frames
for _, frame := range frames {
    if frame.IsRuntime() {
        continue  // Skip runtime frames
    }
    fmt.Printf("File: %s, Line: %d, Function: %s\n", 
        frame.File, frame.Line, frame.Name)
}
```

### Caller Identification

Identify specific callers with configurable skip levels:

```go
// Caller(0) - current function
// Caller(1) - caller of current function
// Caller(2) - caller's caller, etc.

func example() {
    // This will return information about whoever called example()
    caller := stack.Caller(1)
    fmt.Printf("Called from %s:%d\n", caller.File, caller.Line)
}
```

## Advanced Usage

### Selective Stack Capture

```go
import "github.com/pubgo/funk/v2/stack"

// Capture specific number of frames
frames := stack.Callers(10)  // Capture up to 10 frames

// Capture frames with skip
frames := stack.Callers(5, 2)  // Skip 2 frames, then capture 5
```

### Stack Frame Filtering

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/samber/lo"
)

// Filter out runtime frames
frames := lo.Filter(stack.Trace(), func(frame *stack.Frame, _ int) bool {
    return !frame.IsRuntime()
})

// Filter by package
appFrames := lo.Filter(frames, func(frame *stack.Frame, _ int) bool {
    return strings.HasPrefix(frame.Pkg, "github.com/mycompany/myapp")
})
```

### Frame Formatting

```go
import "github.com/pubgo/funk/v2/stack"

frame := stack.Caller(0)

// Full format
fmt.Printf("Full: %s\n", frame.String())  // /path/to/file.go:123 FunctionName

// Short format
fmt.Printf("Short: %s\n", frame.Short())  // file.go:123 FunctionName
```

### Performance Optimization

The module uses caching to optimize repeated stack operations:

```go
// Repeated calls to stack functions benefit from internal caching
for i := 0; i < 1000; i++ {
    caller := stack.Caller(0)  // Cached after first call
    // Process caller
}
```

## API Reference

### Core Functions

| Function | Description |
|----------|-------------|
| `Trace()` | Capture complete stack trace |
| `Caller(skip int)` | Get caller at specified skip level |
| `Callers(depth int, skips ...int)` | Get multiple callers |
| `CallerWithFunc(fn any)` | Get caller information for function |
| `Stack(p uintptr)` | Get frame for program counter |

### Frame Methods

| Method | Description |
|--------|-------------|
| `String()` | Full frame representation |
| `Short()` | Shortened frame representation |
| `IsRuntime()` | Check if frame is in Go runtime |

### Utility Functions

| Function | Description |
|----------|-------------|
| `GetGORoot()` | Get Go runtime root path |
| `GetStack(skip int)` | Get program counter for stack position |

## Best Practices

1. **Use Appropriate Skip Levels**: Choose correct skip levels for accurate caller identification
2. **Filter Runtime Frames**: Remove runtime frames for cleaner stack traces
3. **Cache When Possible**: Leverage built-in caching for repeated operations
4. **Limit Depth**: Use appropriate stack depth to avoid performance issues
5. **Contextual Information**: Include relevant contextual information in error reports

## Integration Patterns

### With Error Handling

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/pubgo/funk/v2/errors"
)

func wrapWithErrorContext() error {
    caller := stack.Caller(1)
    err := someOperation()
    if err != nil {
        // Wrap error with caller context
        return errors.Wrapf(err, "failed at %s:%d in %s", 
            caller.File, caller.Line, caller.Name)
    }
    return nil
}
```

### With Logging

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/pubgo/funk/v2/log"
)

func logWithContext() {
    caller := stack.Caller(0)
    logger := log.GetLogger("stack")
    
    logger.Info().
        Str("file", caller.File).
        Int("line", caller.Line).
        Str("function", caller.Name).
        Msg("Operation performed")
}
```

### With Debugging Tools

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/pubgo/funk/v2/log"
)

func debugFunction() {
    // Capture stack for debugging
    frames := stack.Trace()
    
    logger := log.GetLogger("debug")
    for i, frame := range frames {
        if i > 10 {  // Limit output
            break
        }
        if frame.IsRuntime() {
            continue
        }
        logger.Debug().Msgf("Stack frame %d: %s", i, frame.String())
    }
}
```