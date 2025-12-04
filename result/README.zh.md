# 结果模块

结果模块通过`Result[T]`和`Error`类型引入函数式编程概念到Go错误处理中。它消除了显式的nil检查，提供了安全的组合错误处理和方法链接。

## 功能特性

- **类型安全的错误处理**: 消除nil检查bug的泛型结果类型
- **方法链接**: 用于转换的流畅API
- **模式匹配**: 使用Match函数进行详尽的错误处理
- **异步支持**: Future类型用于异步计算
- **标准兼容性**: 与现有Go错误处理模式协作

## 安装

```bash
go get github.com/pubgo/funk/v2/result
```

## 快速开始

### 基本用法

```go
import "github.com/pubgo/funk/v2/result"

// 创建成功结果
success := result.OK(42)

// 创建失败结果
failure := result.Fail[int](errors.New("计算失败"))

// 安全解包
if value, ok := success.TryUnwrap(); ok {
    fmt.Printf("值: %d\n", value)
}

// 方法链接
processed := result.OK("hello").
    Map(strings.ToUpper).
    FlatMap(func(s string) result.Result[string] {
        if len(s) > 0 {
            return result.OK(s + " world")
        }
        return result.Fail[string](errors.New("空字符串"))
    })
```

### 错误处理

```go
// 模式匹配
result.OK(42).Match(
    func(value int) { fmt.Printf("成功: %d\n", value) },
    func(err error) { fmt.Printf("错误: %v\n", err) }
)

// 仅错误操作
result.ErrOf(errors.New("出错了")).
    Log().
    InspectErr(func(err error) { 
        // 额外的错误处理
    })
```

## 核心概念

### Result[T]

`Result[T]`表示类型`T`的成功值或错误。它提供了安全的替代方案来返回`(T, error)`对。

#### 创建函数

- `OK(v T)`: 创建成功结果
- `Fail[T](err error)`: 创建失败结果
- `Wrap(v T, err error)`: 从值/错误对创建结果
- `WrapFn(fn func() (T, error))`: 从函数创建结果

#### 转换方法

- `Map(func(T) T)`: 转换成功值
- `FlatMap(func(T) Result[T])`: 转换潜在的新错误
- `Validate(func(T) error)`: 验证带潜在错误的值
- `Or(default T)`: 为错误提供后备值

#### 消费方法

- `Unwrap()`: 获取值或在错误时panic
- `TryUnwrap()`: 安全获取带成功指示的值
- `UnwrapOr(default T)`: 获取值或在错误时使用默认值
- `Expect(msg string)`: 获取值或使用自定义消息panic

### Error

`Error`是专用于仅错误操作的结果类型，提供错误处理和日志记录的实用程序。

#### 创建函数

- `ErrOf(err error)`: 从错误创建
- `ErrOfFn(fn func() error)`: 从错误返回函数创建
- `Errorf(format string, args ...any)`: 创建格式化错误

#### 错误处理方法

- `Log()`: 记录带堆栈跟踪的错误
- `InspectErr(func(error))`: 处理错误而不消费它
- `Match(func(), func(error))`: 模式匹配错误存在

### Future[T]

`Future[T]`表示将最终产生`Result[T]`的异步计算。

#### 创建函数

- `Async(fn func() Result[T])`: 从异步函数创建
- `AsyncErr(fn func() Error)`: 从异步错误函数创建

#### 使用

```go
// 开始异步计算
future := result.Async(func() result.Result[int] {
    // 模拟工作
    time.Sleep(time.Second)
    return result.OK(42)
})

// 等待结果
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result := future.Await(ctx)
```

## 高级用法

### 链式操作

```go
// 复杂转换管道
result := getUser(id).
    FlatMap(validateUser).
    Map(encryptUserData).
    FlatMap(saveUser).
    Inspect(func(savedUser User) {
        log.Info().Msgf("用户保存: %s", savedUser.ID)
    }).
    MapErr(func(err error) error {
        return fmt.Errorf("处理用户失败: %w", err)
    })
```

### 处理集合

```go
// 处理多个结果
results := []result.Result[int]{
    result.OK(1),
    result.OK(2),
    result.Fail[int](errors.New("转换失败")),
    result.OK(4),
}

// 收集所有成功值或第一个错误
collected := result.Collect(results)
values, err := collected.UnwrapErr()

// 分区为成功和失败
successes, failures := result.Partition(results)
```

### 模式匹配

```go
// 详尽的错误处理
result.OK(42).Match(
    func(value int) {
        fmt.Printf("成功: %d\n", value)
    },
    func(err error) {
        fmt.Printf("错误: %v\n", err)
    }
)

// 产生结果的模式匹配
result.OK(42).MatchWithResult(
    func(value int) result.Result[string] {
        return result.OK(strconv.Itoa(value))
    },
    func(err error) result.Result[string] {
        return result.Fail[string](err)
    }
)
```

## 集成模式

### 与标准错误处理

```go
// 从标准(T, error)模式转换
func divide(a, b float64) result.Result[float64] {
    if b == 0 {
        return result.Fail[float64](errors.New("除零"))
    }
    return result.OK(a / b)
}

// 转换为标准(T, error)模式
value, err := divide(10, 2).UnwrapErr()
```

### 与上下文

```go
// 上下文感知操作
func fetchUser(ctx context.Context, id string) result.Result[User] {
    user, err := database.GetUser(ctx, id)
    return result.Wrap(user, err)
}

// 带上下文的异步
future := result.Async(func() result.Result[User] {
    return fetchUser(ctx, userID)
})

result := future.Await(ctx)
```

## 最佳实践

1. **优先使用Result[T]而非(T, error)**: 为了更好的组合性和更少的nil检查
2. **使用TryUnwrap进行安全提取**: 在不确定的情况下避免panic
3. **链式操作以提高可读性**: 使用Map/FlatMap进行数据转换管道
4. **适当处理错误**: 使用Match进行详尽的错误处理
5. **利用Async进行IO绑定操作**: 使用Future[T]进行非阻塞计算

## API参考

### Result[T] 创建

| 函数 | 描述 |
|------|------|
| `OK(v T)` | 创建成功结果 |
| `Fail[T](err error)` | 创建失败结果 |
| `Wrap(v T, err error)` | 从值/错误对创建 |
| `WrapFn(fn func() (T, error))` | 从函数创建 |

### Result[T] 检查

| 方法 | 描述 |
|------|------|
| `IsOK() bool` | 检查结果是否成功 |
| `IsErr() bool` | 检查结果是否错误 |
| `TryUnwrap() (T, bool)` | 安全值提取 |

### Result[T] 转换

| 方法 | 描述 |
|------|------|
| `Map(func(T) T)` | 转换成功值 |
| `FlatMap(func(T) Result[T])` | 转换带潜在新错误 |
| `Validate(func(T) error)` | 验证带潜在错误 |

### Result[T] 消费

| 方法 | 描述 |
|------|------|
| `Unwrap() T` | 获取值或panic |
| `UnwrapOr(default T) T` | 获取值或默认值 |
| `Expect(msg string) T` | 获取值或带消息panic |

### Error 创建

| 函数 | 描述 |
|------|------|
| `ErrOf(err error)` | 从错误创建 |
| `ErrOfFn(fn func() error)` | 从错误函数创建 |
| `Errorf(format string, args ...any)` | 创建格式化错误 |

### 错误处理

| 方法 | 描述 |
|------|------|
| `Log()` | 记录带上下文的错误 |
| `InspectErr(func(error))` | 处理错误而不消费 |
| `Match(func(), func(error))` | 模式匹配错误存在 |