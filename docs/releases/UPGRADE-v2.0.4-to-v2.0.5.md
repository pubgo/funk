# 从 v2.0.4 升级到 v2.0.5

```bash
go get github.com/pubgo/funk/v2@v2.0.5
```

## 推荐迁移模式

### 1. 日志里用完整错误链

**以前：**

```go
log.Error().Err(err).Msg(err.Error()) // 只有根消息
```

**现在：**

```go
log.Error().Err(err).Msg("operation failed")
// 自动附加 error_chain、error_tags、error_id、error_detail
```

也可显式读取：

```go
import "github.com/pubgo/funk/v2/errors"

log.Error().
    Str("chain", errors.FormatChain(err)).
    Interface("tags", errors.CollectUserTags(err)).
    Err(err).
    Send()
```

### 2. result 错误展示

```go
r := result.ErrOf(err)
if r.IsErr() {
    fmt.Println(r.Message()) // 多层错误链
    fmt.Println(r.Tags())    // 业务 tags
}
```

### 3. 读取错误元数据

```go
tags := errors.GetTags(err)           // 最内层 *Err 的 tags
all := errors.CollectTags(err)        // 合并整条链（含 wrap msg）
user := errors.CollectUserTags(err)   // 排除系统 wrap msg
```

### 4. errcode 注册与查找

```go
// 推荐
if err := errcode.RegisterErrCode(code); err != nil { ... }

// init 中快速注册
errcode.MustRegisterErrCode(code)

// 按名称查找
code, ok := errcode.LookupErrCode("demo.user.not_found")
```

`RegisterErrCodes` 仍可用（内部 panic），生成代码无需改动。

## 行为变化检查清单

| 场景 | v2.0.4 | v2.0.5 |
|------|--------|--------|
| `Wrap` 后 `Error()` | 根消息 | 不变 |
| 完整可读链 | 手动拼接 | `FormatChain` / `log.Err` |
| `WrapStack` | 额外打 stderr | 仅采集栈 |
| `JsonPrint` 失败 | panic | 返回 `nil` |
| `NewCodeErrWithMsg` | `ToTitle` | 保留原大小写 |

## 文档

- [v2.0.5 发布说明](./v2.0.5.md) / [中文版](./v2.0.5.zh.md)
- [errors/README.md](../../errors/README.md)
- [errors/errcode/README.md](../../errors/errcode/README.md)
