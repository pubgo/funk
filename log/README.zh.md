# 日志模块

日志模块提供了一个基于zerolog的高性能、结构化日志系统，支持上下文感知日志、错误详细信息捕获和模块化日志记录器管理。

## 功能特性

- **高性能日志记录**: 基于zerolog，最小开销
- **结构化日志**: JSON格式日志，自动字段提取
- **上下文感知日志**: 请求范围日志，自动字段注入
- **错误详细信息捕获**: 自动错误丰富带堆栈跟踪和元数据
- **模块化日志记录器**: 带命名空间的不同组件日志记录器
- **多种输出格式**: 控制台和JSON输出格式

## 设计说明

- 该包当前**有意**使用 `zerolog` 作为具体后端，并通过 `Event` 别名暴露 `zerolog.Event`，以保留低开销和流式链式 API。
- `Logger.WithFields(...)` 会返回一个派生 logger，不会修改原 logger。
- `WithLogger(ctx, ...)` 是给 `context.Context` 绑定 logger 的便捷入口。
- `WithFields(ctx, ...)` 是给 `context.Context` 添加请求级字段的便捷入口。
- `CreateFieldsCtx(...)` 会对初始字段 map 做一次快照，避免调用方后续修改原 map 时污染日志上下文。
- `WithLogger(...)`、`WithFields(...)`、`GetFieldsFromCtx(...)` 和 `WithDisabled(...)` 都是 nil-safe 的辅助函数，适合请求级日志链路使用。
- `log.Info(ctx)`、`log.Err(err, ctx)` 这类全局 helper 会自动尊重 context 中绑定的 logger。
- `FromCtx(ctx)` 是从 context 中取出“实际生效 logger”的更顺手入口，适合继续链式调用。
- `UpdateFieldsCtx(...)` 会返回一个新的 context，不会修改父 context 中已有的字段 map。
- 当 logger 默认字段与 context 字段同名时，**context 字段优先**，这样请求级元数据可以安全覆盖模块默认值。

## 推荐用法

- **常规请求链路**：优先使用 `WithLogger` + `WithFields`，然后直接调用 `log.Info(ctx)` / `log.Err(err, ctx)`。
- **需要继续链式操作 logger**：使用 `log.FromCtx(ctx)`。
- **底层显式构造场景**：保留 `CreateCtx` / `CreateFieldsCtx`，用于更严格、非 helper 风格的用法。

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
ctx = log.WithLogger(ctx, log.GetLogger("request"))
ctx = log.WithFields(ctx, map[string]any{
    "request_id": "abc123",
    "user_id": 42,
})

// 带上下文的日志记录（字段自动包含）
log.Info(ctx).Msg("处理请求")

// 全局 helper 也会自动使用 context 里绑定的 logger。

// 如果还需要继续链式调用 logger，可从 context 中取出实际生效 logger。
log.FromCtx(ctx).WithName("db").Info(ctx).Msg("查询完成")
```

### Context Helper 速查

- `WithLogger`：给 context 绑定请求级 logger
- `WithFields`：追加请求级字段
- `log.Info(ctx)` / `log.Err(err, ctx)`：最常用的日志输出路径
- `FromCtx`：当你还需要继续链式调用 logger 时使用
- `CreateCtx` / `CreateFieldsCtx`：更底层、显式的构造入口

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

### Slog / Std 适配器用法

```go
import (
    "log/slog"
    "github.com/pubgo/funk/v2/log"
)

ctx := log.WithLogger(nil, log.GetLogger("request"))

// slog 适配器
slogger := slog.New(log.NewSlog(log.FromCtx(ctx)))
slogger.Info("handled request")

// std 风格适配器
std := log.NewStd(log.FromCtx(ctx))
std.Println("request", "finished")
```

说明：

- `NewSlog(nil)` 和 `NewStd(nil)` 都会安全回退到全局 logger。
- `StdLogger.Println(...)` 现在会像标准库一样，在参数之间插入空格。

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

| 函数                                                  | 描述                                  |
| ----------------------------------------------------- | ------------------------------------- |
| `FromCtx(ctx context.Context)`                        | 获取 context 中实际生效的 logger      |
| `GetFromCtx(ctx context.Context, loggers ...Logger)`  | 底层 logger 查找入口，可传备用 logger |
| `GetLogger(names ...string)`                          | 获取模块特定日志记录器                |
| `SetLogger(log *zerolog.Logger)`                      | 设置全局日志记录器                    |
| `CreateCtx(ctx context.Context, ll Logger)`           | 创建带显式 logger 的上下文            |
| `WithLogger(ctx context.Context, ll Logger)`          | 给上下文绑定 logger                   |
| `WithFields(ctx context.Context, fields Fields)`      | 向上下文添加字段                      |
| `CreateFieldsCtx(ctx context.Context, evt Fields)`    | 创建带初始字段的上下文                |
| `UpdateFieldsCtx(ctx context.Context, fields Fields)` | 在派生 context 中新增或覆盖字段       |
| `GetFieldsFromCtx(ctx context.Context)`               | 从上下文提取字段                      |

### 日志函数

| 函数                                     | 描述                             |
| ---------------------------------------- | -------------------------------- |
| `Debug(ctx ...context.Context)`          | 开始调试级别日志事件             |
| `Info(ctx ...context.Context)`           | 开始信息级别日志事件             |
| `Warn(ctx ...context.Context)`           | 开始警告级别日志事件             |
| `Error(ctx ...context.Context)`          | 开始错误级别日志事件             |
| `Err(err error, ctx ...context.Context)` | 开始带错误详细信息的错误日志事件 |
| `Fatal(ctx ...context.Context)`          | 开始致命日志事件                 |
| `Panic(ctx ...context.Context)`          | 开始恐慌日志事件                 |

### 实用函数

| 函数                                        | 描述                                    |
| ------------------------------------------- | --------------------------------------- |
| `NewEvent()`                                | 创建新字典事件                          |
| `RecordErr(logs ...Logger)`                 | 创建错误记录函数                        |
| `Output(w io.Writer)`                       | 创建带自定义输出的日志记录器            |
| `OutputWriter(w func([]byte) (int, error))` | 创建带自定义写入函数的日志记录器        |
| `NewSlog(log Logger)`                       | 创建基于 `log.Logger` 的 `slog.Handler` |
| `NewStd(log Logger)`                        | 创建 std 风格 logger 适配器             |

## 最佳实践

1. **使用模块日志记录器**: 为更好的组织创建模块特定日志记录器
2. **包含上下文**: 使用上下文感知日志进行请求跟踪
3. **结构化日志**: 使用相关字段进行结构化日志记录
4. **适当级别**: 为不同场景使用适当的日志级别
5. **避免敏感数据**: 永远不要记录密码、令牌或其他敏感信息
6. **错误详细信息**: 始终记录带足够上下文的错误以便调试
7. **把 Logger 当成不可变对象**: 通过 `WithName` / `WithFields` 链式派生新 logger，而不是依赖原地修改
8. **请求数据优先放进 Context**: 请求级字段放入 `context.Context`，模块级默认字段放在 logger 上

## 性能说明

- 该模块保留 `zerolog` 作为执行引擎，以维持它的低分配特性和顺滑的 builder 体验。
- 新增回归测试确保 context 字段合并不会污染可复用 logger 的内部状态。
- `z_bench_test.go` 中提供了 context 合并和 context-aware info logging 的 benchmark，便于后续改动快速比对性能。

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
    ctx = log.WithLogger(ctx, log.GetLogger("request"))
    ctx = log.WithFields(ctx, map[string]any{
        "request_id": requestID,
    })
    
    // 带请求上下文的日志记录
    log.Info(ctx).Msg("处理请求")

    // 需要继续链式调用 logger 时，可从 context 中取出实际生效 logger。
    log.FromCtx(ctx).WithName("downstream").Info(ctx).Msg("调用 processData")
    
    // 将上下文传递给下游函数
    processData(ctx)
}
```