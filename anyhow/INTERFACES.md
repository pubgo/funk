# Interface Abstractions in Anyhow

这个文档描述了 anyhow 库中的接口抽象系统，它允许你编写通用的、类型安全的代码来处理不同的容器类型。

## 概述

anyhow 库提供了一套完整的接口抽象，让你可以编写多态的函数来处理 `Result[T]` 和 `Error` 类型。这些接口遵循函数式编程的原则，提供了一致的 API 来处理可能失败的操作。

> **注意**: 我们移除了 `Option[T]` 类型，因为它在 Go 语言中不够惯用。Go 有零值概念和指针等更自然的方式来处理可选值。

## 核心接口

### 状态检查接口

#### `Checkable`
定义可以检查成功/失败状态的类型：
```go
type Checkable interface {
    IsOk() bool
    IsError() bool
    String() string
}
```

**实现类型**: `Result[T]`, `Error`

**使用示例**:
```go
func CheckStatus(c Checkable) string {
    if c.IsOk() {
        return "✓ Success"
    }
    return "✗ Failed: " + c.String()
}

// 可以用于不同类型
result := Ok(42)
err := ErrorOf(fmt.Errorf("failed"))
fmt.Println(CheckStatus(result)) // ✓ Success
fmt.Println(CheckStatus(err))    // ✗ Failed: Error(failed)
```

#### `Unwrappable[T]`
定义可以安全解包值的类型：
```go
type Unwrappable[T any] interface {
    Unwrap() T
    UnwrapOr(defaultValue T) T
}
```

**实现类型**: `Result[T]`

**使用示例**:
```go
func SafeUnwrap[T any](u Unwrappable[T], defaultValue T) T {
    return u.UnwrapOr(defaultValue)
}

result := Ok(42)
value := SafeUnwrap(result, 0) // 42

// For optional values, use Go's natural patterns:
var optionalString string = "hello"  // or "" for no value
if optionalString != "" {
    // use optionalString
}
```

#### `ErrorAccessible`
定义可以提供错误信息的类型：
```go
type ErrorAccessible interface {
    Err() error
}
```

**实现类型**: `Result[T]`, `Error`

### 函数式编程接口

#### `Mappable[T]`
定义支持值转换的类型：
```go
type Mappable[T any] interface {
    Map(fn func(T) T) Mappable[T]
}
```

#### `Flatmappable[T]`
定义支持单子绑定操作的类型：
```go
type Flatmappable[T any] interface {
    FlatMap(fn func(T) Flatmappable[T]) Flatmappable[T]
}
```

#### `Filterable[T]`
定义支持条件过滤的类型：
```go
type Filterable[T any] interface {
    Filter(predicate func(T) bool, errorMsg string) Filterable[T]
}
```

#### `Inspectable[T]`
定义支持副作用操作的类型：
```go
type Inspectable[T any] interface {
    Inspect(fn func(T)) Inspectable[T]
}
```

### 错误处理接口

#### `Recoverable[T]`
定义支持错误恢复的类型：
```go
type Recoverable[T any] interface {
    OrElse(fn func(error) Recoverable[T]) Recoverable[T]
}
```

#### `ErrorMappable`
定义支持错误转换的类型：
```go
type ErrorMappable interface {
    MapErr(fn func(error) error) ErrorMappable
}
```

#### `Panicable[T]`
定义可以在错误时 panic 的类型：
```go
type Panicable[T any] interface {
    Must() T
    Expect(msg string) T
}
```

## 组合接口

### `Monad[T]`
组合基本的单子操作：
```go
type Monad[T any] interface {
    Mappable[T]
    Flatmappable[T]
    Inspectable[T]
}
```

### `ErrorHandlingMonad[T]`
扩展单子以支持错误处理：
```go
type ErrorHandlingMonad[T any] interface {
    Monad[T]
    ErrorAccessible
    ErrorMappable
    ErrorInspectable
    Recoverable[T]
    Must() T
    Expect(msg string) T
}
```

### `Container[T]`
表示通用容器：
```go
type Container[T any] interface {
    Checkable
    Unwrappable[T]
    Inspectable[T]
}
```

## 实用工具

### 错误收集器

`ErrorCollector` 可以从多个 `ErrorAccessible` 类型收集错误：

```go
collector := NewErrorCollector()

result := Fail[string](fmt.Errorf("result error"))
err := ErrorOf(fmt.Errorf("error type error"))
okResult := Ok("success")

collector.Collect(result).Collect(err).Collect(okResult)

if collector.HasErrors() {
    fmt.Printf("Collected %d errors\n", len(collector.Errors()))
    joinedErr := collector.JoinedError()
    // 处理合并的错误
}
```

### 批量操作

#### `BatchUnwrap`
批量解包多个 `Unwrappable` 值：
```go
items := []Unwrappable[int]{
    Ok(1),
    Ok(2),
    Ok(3),
}

results := BatchUnwrap(items, 0) // [1, 2, 3]
```

#### `ProcessCheckable`
处理多个 `Checkable` 项目：
```go
items := []Checkable{
    Ok(42),
    Fail[string](fmt.Errorf("error")),
    ErrorOf(nil),
    ErrorOf(fmt.Errorf("another error")),
}

successes, failures := ProcessCheckable(items...) // 2, 2
```

#### `MapUnwrappable`
对 `Unwrappable` 值应用转换：
```go
items := []Unwrappable[int]{Ok(1), Ok(2), Ok(3)}
doubled := MapUnwrappable(items, func(x int) int { return x * 2 }, 0)
// [2, 4, 6]
```

### 验证管道

`ValidationPipeline` 允许你组合多个验证步骤：

```go
// 定义验证步骤
type LengthValidator struct {
    MinLength int
}

func (lv LengthValidator) Validate(s string) Checkable {
    if len(s) >= lv.MinLength {
        return WrapBoolAsCheckable(true, "length valid")
    }
    return WrapBoolAsCheckable(false, "length too short")
}

// 创建管道
pipeline := NewValidationPipeline[string]().
    AddStep(LengthValidator{MinLength: 3}).
    AddStep(ContainsValidator{Required: "a"})

// 使用管道
if pipeline.IsValid("abc") {
    fmt.Println("Valid input")
}

results := pipeline.Validate("ab")
for _, result := range results {
    if result.IsError() {
        fmt.Println("Validation failed:", result.String())
    }
}
```

### 适配器

#### `WrapBoolAsCheckable`
将布尔值包装为 `Checkable` 接口：
```go
checkable := WrapBoolAsCheckable(true, "success")
if checkable.IsOk() {
    fmt.Println("Operation succeeded")
}
```

## 高级用法

### 多态函数

使用接口可以编写适用于多种类型的函数：

```go
// 通用状态检查
func LogStatus[T Checkable](items []T, logger *log.Logger) {
    for i, item := range items {
        if item.IsOk() {
            logger.Printf("Item %d: Success", i)
        } else {
            logger.Printf("Item %d: Failed - %s", i, item.String())
        }
    }
}

// 通用错误收集
func CollectAllErrors(items []ErrorAccessible) []error {
    var errors []error
    for _, item := range items {
        if err := item.Err(); err != nil {
            errors = append(errors, err)
        }
    }
    return errors
}

// 通用值提取
func ExtractValues[T any](items []Unwrappable[T], defaultValue T) []T {
    results := make([]T, len(items))
    for i, item := range items {
        results[i] = item.UnwrapOr(defaultValue)
    }
    return results
}
```

### 函数式组合

```go
// 链式操作
func ProcessData[T any](
    data T,
    validators []func(T) Checkable,
    transformer func(T) T,
) Result[T] {
    // 验证数据
    for _, validate := range validators {
        if validate(data).IsError() {
            return Fail[T](fmt.Errorf("validation failed"))
        }
    }
    
    // 转换数据
    return Ok(transformer(data))
}
```

## 设计原则

1. **类型安全**: 所有接口都是类型安全的，编译时检查
2. **组合性**: 小接口可以组合成更大的接口
3. **一致性**: 所有类型都遵循相同的命名约定
4. **可扩展性**: 用户可以实现这些接口来创建自己的类型
5. **向后兼容**: 不破坏现有 API

## 性能考虑

- 接口调用有轻微的运行时开销
- 对于性能关键的代码，考虑直接使用具体类型
- 批量操作函数经过优化，适合处理大量数据

## 最佳实践

1. **优先使用小接口**: 只依赖你需要的方法
2. **组合而非继承**: 使用接口组合来构建复杂行为
3. **明确错误处理**: 使用 `ErrorAccessible` 来处理错误
4. **利用类型推断**: Go 的类型推断可以简化泛型代码
5. **文档化约束**: 清楚地记录接口的预期行为

## 示例：完整的数据处理管道

```go
func ProcessUserData(rawData string) Result[User] {
    return Try(func() (User, error) {
        // 解析 JSON
        var userData map[string]interface{}
        if err := json.Unmarshal([]byte(rawData), &userData); err != nil {
            return User{}, err
        }
        
        // 验证数据
        pipeline := NewValidationPipeline[map[string]interface{}]().
            AddStep(RequiredFieldValidator{Field: "name"}).
            AddStep(RequiredFieldValidator{Field: "email"})
        
        if !pipeline.IsValid(userData) {
            return User{}, fmt.Errorf("validation failed")
        }
        
        // 创建用户
        return User{
            Name:  userData["name"].(string),
            Email: userData["email"].(string),
        }, nil
    }).
    Inspect(func(user User) {
        log.Printf("Created user: %+v", user)
    }).
    MapErr(func(err error) error {
        return fmt.Errorf("failed to process user data: %w", err)
    })
}
```

这个接口系统让 anyhow 库更加灵活和强大，同时保持了类型安全和性能。 