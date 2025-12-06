# 堆栈模块

堆栈模块提供用于堆栈跟踪分析和调用者识别的实用程序，实现详细的调试和错误报告功能。

## 功能特性

- **运行时堆栈跟踪捕获**: 高效捕获运行时堆栈跟踪
- **调用者识别**: 准确识别调用函数
- **堆栈帧分析**: 带元数据的详细堆栈帧分析
- **性能优化**: 缓存堆栈跟踪处理以最小化开销
- **Go运行时集成**: 与Go运行时内部的深度集成

## 安装

```bash
go get github.com/pubgo/funk/v2/stack
```

## 快速开始

### 基本堆栈跟踪捕获

```go
import "github.com/pubgo/funk/v2/stack"

// 捕获当前堆栈跟踪
frames := stack.Trace()

// 打印堆栈帧
for _, frame := range frames {
    fmt.Printf("%s\n", frame.String())
}
```

### 调用者识别

```go
import "github.com/pubgo/funk/v2/stack"

// 获取直接调用者
caller := stack.Caller(0)
fmt.Printf("被调用于: %s:%d %s\n", caller.File, caller.Line, caller.Name)

// 获取父调用者
parent := stack.Caller(1)
fmt.Printf("父调用者: %s:%d %s\n", parent.File, parent.Line, parent.Name)
```

### 基于函数的调用者识别

```go
import "github.com/pubgo/funk/v2/stack"

func myFunction() {
    // 获取此函数的调用者
    caller := stack.CallerWithFunc(myFunction)
    fmt.Printf("myFunction被调用于: %s\n", caller.String())
}
```

## 核心概念

### 堆栈帧

模块使用详细元数据表示堆栈帧：

```go
type Frame struct {
    Name string // 函数名称
    Pkg  string // 包路径
    File string // 文件路径
    Line int    // 行号
}
```

### 堆栈跟踪捕获

捕获完整堆栈跟踪用于调试：

```go
// 捕获完整堆栈跟踪
frames := stack.Trace()

// 处理帧
for _, frame := range frames {
    if frame.IsRuntime() {
        continue  // 跳过运行时帧
    }
    fmt.Printf("文件: %s, 行: %d, 函数: %s\n", 
        frame.File, frame.Line, frame.Name)
}
```

### 调用者识别

识别特定调用者带可配置跳过级别：

```go
// Caller(0) - 当前函数
// Caller(1) - 当前函数的调用者
// Caller(2) - 调用者的调用者，等等

func example() {
    // 这将返回调用example()的信息
    caller := stack.Caller(1)
    fmt.Printf("被调用于 %s:%d\n", caller.File, caller.Line)
}
```

## 高级用法

### 选择性堆栈捕获

```go
import "github.com/pubgo/funk/v2/stack"

// 捕获特定数量的帧
frames := stack.Callers(10)  // 捕获最多10帧

// 带跳过的帧捕获
frames := stack.Callers(5, 2)  // 跳过2帧，然后捕获5帧
```

### 堆栈帧过滤

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/samber/lo"
)

// 过滤掉运行时帧
frames := lo.Filter(stack.Trace(), func(frame *stack.Frame, _ int) bool {
    return !frame.IsRuntime()
})

// 按包过滤
appFrames := lo.Filter(frames, func(frame *stack.Frame, _ int) bool {
    return strings.HasPrefix(frame.Pkg, "github.com/mycompany/myapp")
})
```

### 帧格式化

```go
import "github.com/pubgo/funk/v2/stack"

frame := stack.Caller(0)

// 完整格式
fmt.Printf("完整: %s\n", frame.String())  // /path/to/file.go:123 FunctionName

// 短格式
fmt.Printf("短: %s\n", frame.Short())  // file.go:123 FunctionName
```

### 性能优化

模块使用缓存优化重复堆栈操作：

```go
// 重复调用堆栈函数受益于内部缓存
for i := 0; i < 1000; i++ {
    caller := stack.Caller(0)  // 第一次调用后缓存
    // 处理调用者
}
```

## API参考

### 核心函数

| 函数 | 描述 |
|------|------|
| `Trace()` | 捕获完整堆栈跟踪 |
| `Caller(skip int)` | 获取指定跳过级别的调用者 |
| `Callers(depth int, skips ...int)` | 获取多个调用者 |
| `CallerWithFunc(fn any)` | 获取函数的调用者信息 |
| `Stack(p uintptr)` | 获取程序计数器的帧 |

### 帧方法

| 方法 | 描述 |
|------|------|
| `String()` | 完整帧表示 |
| `Short()` | 缩短的帧表示 |
| `IsRuntime()` | 检查帧是否在Go运行时中 |

### 实用函数

| 函数 | 描述 |
|------|------|
| `GetGORoot()` | 获取Go运行时根路径 |
| `GetStack(skip int)` | 获取堆栈位置的程序计数器 |

## 最佳实践

1. **使用适当的跳过级别**: 选择正确的跳过级别以准确识别调用者
2. **过滤运行时帧**: 移除运行时帧以获得更清晰的堆栈跟踪
3. **尽可能缓存**: 利用内置缓存进行重复操作
4. **限制深度**: 使用适当的堆栈深度以避免性能问题
5. **上下文信息**: 在错误报告中包含相关上下文信息

## 集成模式

### 与错误处理

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/pubgo/funk/v2/errors"
)

func wrapWithErrorContext() error {
    caller := stack.Caller(1)
    err := someOperation()
    if err != nil {
        // 用调用者上下文包装错误
        return errors.Wrapf(err, "失败于 %s:%d 在 %s", 
            caller.File, caller.Line, caller.Name)
    }
    return nil
}
```

### 与日志记录

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
        Msg("操作执行")
}
```

### 与调试工具

```go
import (
    "github.com/pubgo/funk/v2/stack"
    "github.com/pubgo/funk/v2/log"
)

func debugFunction() {
    // 捕获堆栈用于调试
    frames := stack.Trace()
    
    logger := log.GetLogger("debug")
    for i, frame := range frames {
        if i > 10 {  // 限制输出
            break
        }
        if frame.IsRuntime() {
            continue
        }
        logger.Debug().Msgf("堆栈帧 %d: %s", i, frame.String())
    }
}
```