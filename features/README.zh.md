# 功能模块

功能模块通过功能标志提供集中式应用程序配置管理系统。它支持运行时配置、环境变量集成和CLI标志生成。

## 功能特性

- **集中配置**: 单一注册表管理所有功能标志
- **多来源**: 环境变量、CLI标志和程序化配置
- **类型安全**: 强类型功能值（布尔、字符串、整数、浮点、json）
- **元数据支持**: 用于组织和记录功能的标签
- **CLI集成**: urfave/cli的自动CLI标志生成

## 安装

```bash
go get github.com/pubgo/funk/v2/features
```

## 快速开始

### 定义功能标志

```go
import "github.com/pubgo/funk/v2/features"

// 布尔功能标志
var debugMode = features.Bool("debug", false, "启用调试模式")

// 带元数据的字符串功能标志
var logLevel = features.String("log_level", "info", "日志级别",
    map[string]any{
        "group": "logging",
        "allowed": []string{"debug", "info", "warn", "error"},
    })

// 整数功能标志
var maxRetries = features.Int("max_retries", 3, "最大重试次数")

// JSON功能标志
type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}
var serverConfig = features.Json("server_config", &Config{}, "服务器配置")
```

### 使用功能值

```go
// 访问功能值
if debugMode.Value() {
    log.Println("调试模式已启用")
}

// 程序化更新功能值
err := maxRetries.Set("5")
if err != nil {
    log.Printf("设置max_retries失败: %v", err)
}

// 访问带元数据
flag := features.Lookup("log_level")
if flag != nil {
    group := flag.Tags["group"]
    fmt.Printf("日志级别功能在组中: %v\n", group)
}
```

## 核心概念

### 功能注册表

模块使用中央注册表模式管理所有功能标志：

```go
// 默认注册表（最常用）
var debugMode = features.Bool("debug", false, "启用调试模式")

// 自定义注册表
customFeatures := features.NewFeature()
var customFlag = customFeatures.Bool("custom", false, "自定义功能")
```

### 值类型

模块支持几种带适当解析的值类型：

1. **BoolValue**: 布尔标志，灵活解析（"true", "1", "on", "yes"）
2. **StringValue**: 字符串值
3. **IntValue**: 整数值
4. **FloatValue**: 浮点值
5. **JsonValue**: JSON可序列化值

### 元数据和标签

功能标志可以包含通过标签的元数据：

```go
var feature = features.Bool("feature", false, "描述",
    map[string]any{
        "group": "performance",
        "experimental": true,
        "owner": "team-name",
    })
```

## 高级用法

### 环境变量集成

功能标志自动与环境变量集成：

```go
// 名为"debug"的功能标志将从以下位置读取：
// - FEATURE_DEBUG（自动生成）
// - DEBUG（直接映射）

// 自定义环境变量映射
var custom = features.String("custom", "default", "描述")
// 从FEATURE_CUSTOM读取
```

### CLI标志生成

为urfave/cli应用程序生成CLI标志：

```go
import "github.com/pubgo/funk/v2/features/featureflags"

app := &cli.App{
    Flags: append(
        baseFlags,
        featureflags.GetFlags()..., // 添加所有功能标志作为CLI选项
    ),
    Action: func(c *cli.Context) error {
        // 功能值会自动从CLI更新
        if debugMode.Value() {
            log.Println("通过CLI启用调试模式")
        }
        return nil
    },
}
```

### 功能枚举

遍历所有注册的功能：

```go
features.VisitAll(func(flag *features.Flag) {
    fmt.Printf("功能: %s = %v\n", flag.Name, flag.Value.Get())
    
    // 访问元数据
    if group, ok := flag.Tags["group"]; ok {
        fmt.Printf("  组: %v\n", group)
    }
})
```

### 动态更新

功能值可以在运行时更新：

```go
// 程序化更新
err := debugMode.Set("true")
if err != nil {
    log.Printf("启用调试模式失败: %v", err)
}

// 从字符串更新（对配置加载有用）
values := map[string]string{
    "debug": "true",
    "max_retries": "5",
    "log_level": "debug",
}

for name, value := range values {
    flag := features.Lookup(name)
    if flag != nil {
        err := flag.Value.Set(value)
        if err != nil {
            log.Printf("设置%s失败: %v", name, err)
        }
    }
}
```

## 集成模式

### 与配置文件

与配置文件加载结合：

```go
func loadConfigFromFile(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    
    var config map[string]any
    if err := json.Unmarshal(data, &config); err != nil {
        return err
    }
    
    // 从配置更新功能标志
    features.VisitAll(func(flag *features.Flag) {
        if value, exists := config[flag.Name]; exists {
            strValue := fmt.Sprintf("%v", value)
            if err := flag.Value.Set(strValue); err != nil {
                log.Printf("从配置设置%s失败: %v", flag.Name, err)
            }
        }
    })
    
    return nil
}
```

### 与可观测性

与监控系统集成：

```go
// 导出功能值用于监控
func exportFeatures() map[string]any {
    metrics := make(map[string]any)
    features.VisitAll(func(flag *features.Flag) {
        metrics[flag.Name] = flag.Value.Get()
    })
    return metrics
}

// 定期导出到指标系统
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        metrics := exportFeatures()
        // 发送到指标系统（Prometheus等）
        sendMetrics(metrics)
    }
}()
```

## 最佳实践

1. **在包级别定义功能**: 便于在整个应用程序中访问
2. **使用描述性名称和用途字符串**: 使功能自我记录
3. **应用元数据进行组织**: 使用标签对相关功能进行分组
4. **利用环境变量**: 启用无需代码更改的配置
5. **原子更新功能**: 当更改多个相关功能时
6. **记录实验性功能**: 使用元数据标记不稳定功能
7. **监控功能使用**: 跟踪实际使用的功能

## API参考

### 值创建函数

| 函数 | 描述 |
|------|------|
| `Bool(name, default, usage, tags...)` | 布尔功能 | 
| `String(name, default, usage, tags...)` | 字符串功能 |
| `Int(name, default, usage, tags...)` | 整数功能 |
| `Float(name, default, usage, tags...)` | 浮点功能 |
| `Json[T](name, default, usage, tags...)` | JSON功能 |

### 功能注册表

| 函数 | 描述 |
|------|------|
| `Lookup(name string)` | 按名称查找功能 |
| `VisitAll(func(*Flag))` | 遍历所有功能 |
| `NewFeature()` | 创建自定义注册表 |

### 值方法

| 方法 | 描述 |
|------|------|
| `Value() T` | 获取当前值 |
| `Set(string) error` | 从字符串更新 |
| `String() string` | 字符串表示 |
| `Type() ValueType` | 获取值类型 |
| `Get() any` | 作为interface{}获取 |

### 标志属性

| 属性 | 描述 |
|------|------|
| `Name` | 功能名称 |
| `Usage` | 描述字符串 |
| `Value` | 当前值包装器 |
| `Tags` | 元数据映射 |