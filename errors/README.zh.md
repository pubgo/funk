# 错误模块

错误模块为Go应用程序提供增强的错误处理功能，提供丰富的上下文、堆栈跟踪和元数据，同时保持与标准`error`接口的兼容性。

## 功能特性

- **丰富的错误上下文**: 附加元数据和标签到错误
- **堆栈跟踪**: 自动堆栈跟踪捕获用于调试
- **错误包装**: 使用元数据保留错误链
- **gRPC兼容性**: 与gRPC状态代码集成
- **Panic恢复**: 与 `recovery` 包配合使用，安全处理 panic

## 错误消息语义

`Wrap` / `Wrapf` 会把上下文作为元数据附加，而不会改变 `Error()` 的文本：

```go
err := errors.New("数据库连接失败")
wrapped := errors.Wrap(err, "初始化服务失败")

wrapped.Error() // "数据库连接失败" — 保留根错误消息
// 上下文可通过 CollectTags(wrapped)["msg"] 或 wrapped.(*errors.ErrWrap).String() 获取
```

需要完整结构化表示时，请使用 `String()`、`DebugPrint`、`MarshalError` 或 `FormatChain`。

## 安装

```bash
go get github.com/pubgo/funk/v2/errors
```

## 快速开始

### 创建错误

```go
import "github.com/pubgo/funk/v2/errors"

// 简单错误创建
err := errors.New("数据库连接失败")

// 带元数据的错误
err := errors.New("数据库连接失败", errors.Tags{
    "component": "database",
    "host": "localhost:5432",
})
```

### 错误包装

```go
// 用额外上下文包装
err = errors.Wrap(err, "初始化服务失败")

// 用格式化消息包装
err = errors.Wrapf(err, "初始化服务失败，重试%d次后", retryCount)

// 用堆栈跟踪包装
err = errors.WrapStack(err)
```

### 错误检查

```go
// 标准错误检查
if err != nil {
    // 处理错误
}

// 类型断言
if dbErr, ok := errors.AsA[*DatabaseError](err); ok {
    // 处理数据库特定错误
}

// 错误链遍历
for err != nil {
    fmt.Println(err.Error())
    err = errors.Unwrap(err)
}
```

## 核心概念

### 错误类型

该模块提供两种主要错误类型：

1. **`Err`**: 带消息、详细信息和标签的基础错误实现
2. **`ErrWrap`**: 用于向错误添加上下文和堆栈跟踪的包装器

### 错误创建函数

- `New(msg string, tags ...Tags)`: 创建带可选元数据的新错误
- `Errorf(format string, args ...any)`: 创建格式化错误
- `Wrap(err error, msg string)`: 用额外上下文包装错误
- `Wrapf(err error, format string, args ...any)`: 用格式化上下文包装错误
- `WrapStack(err error)`: 用堆栈跟踪信息包装错误

### 错误检查

- `As(err error, target any)`: 带错误链遍历的类型断言
- `Is(err, target error)`: 检查错误链是否包含特定错误
- `Unwrap(err error)`: 检索底层错误

## 高级用法

### 自定义错误类型

```go
type DatabaseError struct {
    Op string
    Table string
    Err error
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("数据库错误在 %s.%s: %v", e.Op, e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}
```

### 错误元数据

```go
// 向错误附加元数据
err := errors.New("支付处理失败", errors.Tags{
    "transaction_id": "txn_12345",
    "amount": 99.99,
    "currency": "USD",
})

// 从根 *Err 获取元数据
if tags := errors.GetTags(err); tags != nil {
    txnID := tags["transaction_id"]
    // 处理交易ID
}

// 合并整条错误链上的 tags
allTags := errors.CollectTags(err)

// 用户 tags 会排除自动包装上下文 TagKeyMessage ("msg")
userTags := errors.CollectUserTags(err)
```

### 可读错误链

```go
wrapped := errors.Wrap(errors.New("数据库连接失败"), "初始化服务")
fmt.Println(wrapped.Error())             // 数据库连接失败
fmt.Println(errors.FormatChain(wrapped)) // 初始化服务: 数据库连接失败
```

### 堆栈跟踪分析

```go
// 捕获堆栈跟踪
err := errors.WrapStack(errors.New("严重故障"))

// 打印详细错误信息
errors.DebugPrint(err)

// 提取结构化错误数据
data := errors.ErrJsonify(err)
```

## 集成

### 与gRPC

详见 [errcode/README.md](./errcode/README.md)。

```go
import "github.com/pubgo/funk/v2/errors/errcode"

errcode.MustRegisterErrCode(&errorpb.ErrCode{
    Name:       "demo.user.not_found",
    Code:       404,
    StatusCode: errorpb.Code_NotFound,
    Message:    "用户不存在",
})

code, ok := errcode.LookupErrCode("demo.user.not_found")
if ok {
    err := errcode.NewCodeErr(code)
    _ = err
}

status := errcode.ConvertErr2Status(errcode.ParseError(err))
```

### 与日志记录

`log.Err(err)` 会自动附加：

- `error_id`
- `error_chain`（`FormatChain`）
- `error_tags`（`CollectUserTags`）
- `error_detail`（JSON）

```go
import "github.com/pubgo/funk/v2/log"

log.Error().Err(err).Msg("操作失败")
```

### 与 result

```go
import "github.com/pubgo/funk/v2/result"

r := result.ErrOf(err)
if r.IsErr() {
    _ = r.Message() // 完整错误链
    _ = r.Tags()     // 用户 tags
}
```

### 手动结构化日志

```go
log.Error().RawJSON("error", errors.JsonPrint(err)).Send()
```

## 最佳实践

1. **始终包装错误**: 在传递错误到调用堆栈时提供上下文
2. **谨慎使用堆栈跟踪**: 仅在关键路径或调试时使用
3. **附加相关元数据**: 包括有助于调试的信息
4. **保留错误链**: 使用`Unwrap()`维护错误因果关系
5. **在适当级别处理错误**: 不要在没有增加价值的情况下捕获和重新抛出

## API参考

### 错误创建

| 函数 | 描述 |
|------|------|
| `New(msg string, tags ...Tags)` | 创建带可选标签的新错误 |
| `Errorf(format string, args ...any)` | 创建格式化错误 |
| `Wrap(err error, msg string)` | 用上下文包装错误 |
| `Wrapf(err error, format string, args ...any)` | 用格式化上下文包装错误 |
| `WrapStack(err error)` | 用堆栈跟踪包装错误 |
| `WrapTags(err error, tags Tags)` | 用元数据标签包装错误 |

### 错误检查

| 函数 | 描述 |
|------|------|
| `As(err error, target any)` | 带链遍历的类型断言 |
| `Is(err, target error)` | 检查错误链中的特定错误 |
| `Unwrap(err error)` | 获取底层错误 |
| `GetErrorId(err error)` | 获取唯一错误标识符 |
| `GetTags(err error)` | 获取最内层 `*Err` 上的 tags |
| `CollectTags(err error)` | 合并错误链各层的 tags |
| `CollectUserTags(err error)` | 合并用户 tags，并排除包装层 `msg` |
| `RootCause(err error)` | 返回最底层错误 |
| `Walk(err error, fn func(error) bool)` | 遍历错误链 |
| `FormatChain(err error)` | 将错误链拼接为单条消息 |
| `FullMessage(err error)` | `FormatChain` 的别名 |
| `MarshalError(err error)` | 序列化为 JSON 并返回 error |

### 实用函数

| 函数 | 描述 |
|------|------|
| `DebugPrint(err error)` | 带堆栈跟踪的漂亮打印错误 |
| `JsonPrint(err error)` | 序列化为 JSON（失败时返回 nil） |
| `ErrJsonify(err error)` | 将错误转换为结构化数据 |