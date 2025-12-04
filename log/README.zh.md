# 日志模块

日志模块提供了一个基于zerolog的高性能、结构化日志系统，支持上下文感知日志、错误详细信息捕获和模块化日志记录器管理。

## 功能特性

- **高性能日志记录**: 基于zerolog，最小开销
- **结构化日志**: JSON格式日志，自动字段提取
- **上下文感知日志**: 请求范围日志，自动字段注入
- **错误详细信息捕获**: 自动错误丰富带堆栈跟踪和元数据
- **模块化日志记录器**: 带命名空间的不同组件日志记录器
- **多种输出格式**: 控制台和JSON输出格式

## 安装

```bash
go get github.com/pubgo/funk/v2/log
```

## 快速开始

### 基本日志记录

```go
import "github.com/pubgo/funk/v2/log"

// 使用全局日志记录器
log.Info().Msg("应用程序已启动")

// 带字段的日志记录
log.Info().Str("component", "server").Int("port", 8080).Msg("服务器监听中")

// 带自动错误详细信息捕获的错误日志
err := someOperation()
if err != nil {
    log.Err(err).Msg("操作失败")
}
```

### 模块化日志记录

```go
import "github.com/pubgo/funk/v2/log"

// 为特定模块创建日志记录器
logger := log.GetLogger("database")

// 使用模块日志记录器
logger.Info().Msg("连接数据库")
logger.Err(dbError).Msg("数据库连接失败")
```

### 上下文感知日志记录

```go
import (
    "context"
    "github.com/pubgo/funk/v2/log"
)

// 创建带字段的上下文
ctx := context.Background()
ctx = log.WithFields(ctx, map[string]any{
    "request_id": "abc123",
    "user_id": 42,
})

// 带上下文的日志记录（字段自动包含）
log.Info(ctx).Msg("处理请求")
```

## 核心概念

### 日志记录器接口

模块提供`Logger`接口用于一致的日志操作：

```go
type Logger interface {
    Debug(ctx ...context.Context) *zerolog.Event
    Info(ctx ...context.Context) *zerolog.Event
    Warn(ctx ...context.Context) *zerolog.Event
    Error(ctx ...context.Context) *zerolog.Event
    Err(err error, ctx ...context.Context) *zerolog.Event
    // ... 其他方法
}
```

### 全局日志记录器

全局日志记录器可用于简单用例：

```go
// 直接日志函数
log.Debug().Msg("调试消息")
log.Info().Msg("信息消息")
log.Warn().Msg("警告消息")
log.Error().Msg("错误消息")
log.Fatal().Msg("致命消息")
log.Panic().Msg("恐慌消息")
```

### 模块日志记录器

模块特定日志记录器提供命名空间和组件隔离：

```go
// 为特定模块创建日志记录器
dbLogger := log.GetLogger("database")
cacheLogger := log.GetLogger("cache")

// 使用模块特定日志记录器
dbLogger.Info().Msg("数据库操作")
cacheLogger.Info().Msg("缓存操作")
```

### 上下文集成

上下文感知日志记录自动包含相关字段：

```go
// 向上下文添加字段
ctx = log.WithFields(ctx, map[string]any{
    "trace_id": generateTraceID(),
    "span_id": generateSpanID(),
})

// 带上下文的日志记录（字段自动包含）
log.Info(ctx).Msg("带跟踪上下文的处理")
```

## 高级用法

### 自定义日志记录器配置

```go
import (
    "os"
    "github.com/rs/zerolog"
    "github.com/pubgo/funk/v2/log"
)

// 创建自定义zerolog日志记录器
customLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

// 设置为全局日志记录器
log.SetLogger(&customLogger)
```

### 日志记录器级别

```go
import "github.com/rs/zerolog"

// 设置全局日志级别
zerolog.SetGlobalLevel(zerolog.InfoLevel)

// 创建带特定级别的日志记录器
logger := log.GetLogger("module").WithLevel(zerolog.DebugLevel)
```

### 错误详细信息捕获

模块自动丰富错误日志带详细信息：

```go
import "github.com/pubgo/funk/v2/errors"

// 创建带元数据的错误
err := errors.New("数据库连接失败", errors.Tags{
    "host": "localhost",
    "port": 5432,
})

// 记录错误（自动包含错误详细信息）
log.Err(err).Msg("连接尝试失败")
```

### 事件构建

```go
// 构建复杂日志事件
log.Info().
    Str("component", "auth").
    Int("user_id", 12345).
    Dict("metadata", log.NewEvent().
        Str("ip", "192.168.1.1").
        Str("user_agent", "Mozilla/5.0")).
    Msg("用户登录")
```

## API参考

### 核心函数

| 函数 | 描述 |
|------|------|
| `GetLogger(names ...string)` | 获取模块特定日志记录器 |
| `SetLogger(log *zerolog.Logger)` | 设置全局日志记录器 |
| `WithFields(ctx context.Context, fields Fields)` | 向上下文添加字段 |
| `GetFieldsFromCtx(ctx context.Context)` | 从上下文提取字段 |

### 日志函数

| 函数 | 描述 |
|------|------|
| `Debug(ctx ...context.Context)` | 开始调试级别日志事件 |
| `Info(ctx ...context.Context)` | 开始信息级别日志事件 |
| `Warn(ctx ...context.Context)` | 开始警告级别日志事件 |
| `Error(ctx ...context.Context)` | 开始错误级别日志事件 |
| `Err(err error, ctx ...context.Context)` | 开始带错误详细信息的错误日志事件 |
| `Fatal(ctx ...context.Context)` | 开始致命日志事件 |
| `Panic(ctx ...context.Context)` | 开始恐慌日志事件 |

### 实用函数

| 函数 | 描述 |
|------|------|
| `NewEvent()` | 创建新字典事件 |
| `RecordErr(logs ...Logger)` | 创建错误记录函数 |
| `Output(w io.Writer)` | 创建带自定义输出的日志记录器 |
| `OutputWriter(w func([]byte) (int, error))` | 创建带自定义写入函数的日志记录器 |

## 最佳实践

1. **使用模块日志记录器**: 为更好的组织创建模块特定日志记录器
2. **包含上下文**: 使用上下文感知日志进行请求跟踪
3. **结构化日志**: 使用相关字段进行结构化日志记录
4. **适当级别**: 为不同场景使用适当的日志级别
5. **避免敏感数据**: 永远不要记录密码、令牌或其他敏感信息
6. **错误详细信息**: 始终记录带足够上下文的错误以便调试

## 集成模式

### 与错误处理

```go
import (
    "github.com/pubgo/funk/v2/log"
    "github.com/pubgo/funk/v2/errors"
    "github.com/pubgo/funk/v2/result"
)

func processRequest() result.Error {
    err := someOperation()
    if err != nil {
        // 记录带完整详细信息的错误
        log.Err(err).Msg("请求处理失败")
        return result.ErrOf(err)
    }
    return result.Error{}
}
```

### 与配置

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
    Msg("日志配置完成")
```

### 与上下文和跟踪

```go
import (
    "context"
    "github.com/pubgo/funk/v2/log"
    "github.com/google/uuid"
)

func handleRequest(ctx context.Context) {
    // 向上下文添加请求ID
    requestID := uuid.New().String()
    ctx = log.WithFields(ctx, map[string]any{
        "request_id": requestID,
    })
    
    // 带请求上下文的日志记录
    log.Info(ctx).Msg("处理请求")
    
    // 将上下文传递给下游函数
    processData(ctx)
}
```