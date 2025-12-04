# 环境模块

环境模块提供用于处理环境变量的实用程序，包括标准化访问、类型安全检索和.env文件加载。

## 功能特性

- **标准化环境访问**: 环境变量的一致命名约定
- **类型安全检索**: 安全转换为常见类型（布尔、整数、浮点）
- **.env文件支持**: 使用godotenv自动加载.env文件
- **环境扩展**: 使用envsubst的变量替换
- **灵活查找**: 多回退机制的变量解析

## 安装

```bash
go get github.com/pubgo/funk/v2/env
```

## 快速开始

### 基本环境访问

```go
import "github.com/pubgo/funk/v2/env"

// 获取带回退的环境变量
host := env.GetOr("SERVER_HOST", "localhost")

// 获取环境变量或在未设置时panic
port := env.MustGet("SERVER_PORT")

// 使用多个回退名称获取
dbHost := env.Get("DB_HOST", "DATABASE_HOST", "MYSQL_HOST")
```

### 类型安全访问

```go
import "github.com/pubgo/funk/v2/env"

// 布尔值（支持：true, 1, on, yes）
debug := env.GetBool("DEBUG", "ENABLE_DEBUG")

// 整数值
port := env.GetInt("PORT", "SERVER_PORT")

// 浮点值
timeout := env.GetFloat("TIMEOUT", "REQUEST_TIMEOUT")
```

### .env文件加载

```go
import "github.com/pubgo/funk/v2/env"

// 加载.env文件
env.LoadFiles(".env", "config/.env.local")

// 从加载的文件访问变量
apiKey := env.Get("API_KEY")
```

## 核心概念

### 环境变量规范化

模块规范化环境变量名称以确保一致性：

```go
// 这些在规范化后都是等效的：
// server-host => SERVER_HOST
// server.host => SERVER_HOST
// server/host => SERVER_HOST
// SERVER-HOST => SERVER_HOST

host := env.Get("server-host")  // 匹配SERVER_HOST
```

### 多名称查找

函数支持多个名称以实现灵活的变量解析：

```go
// 查找DB_HOST，然后DATABASE_HOST，然后MYSQL_HOST
dbHost := env.Get("DB_HOST", "DATABASE_HOST", "MYSQL_HOST")
```

### 安全类型转换

类型安全获取器提供优雅的转换错误处理：

```go
// GetInt在转换错误时返回-1
port := env.GetInt("PORT")
if port == -1 {
    port = 8080  // 默认值
}

// GetBool处理各种真值
debug := env.GetBool("DEBUG")  // 对"true", "1", "on", "yes"为true
```

## 高级用法

### 环境变量扩展

```go
import "github.com/pubgo/funk/v2/env"

// 在字符串中扩展环境变量
template := "http://${{HOST}}:${{PORT}}"
expanded := env.Expand(template).UnwrapOr(template)
```

### 环境操作

```go
import "github.com/pubgo/funk/v2/env"

// 设置环境变量
env.Set("CUSTOM_VAR", "value").Log()

// 删除环境变量
env.Delete("UNUSED_VAR").Log()

// 检查变量是否存在
if val, exists := env.Lookup("OPTIONAL_VAR"); exists {
    fmt.Printf("可选变量: %s\n", val)
}
```

### 环境映射

```go
import "github.com/pubgo/funk/v2/env"

// 将所有环境变量获取为映射
allVars := env.Map()

// 处理所有变量
for key, value := range allVars {
    fmt.Printf("%s=%s\n", key, value)
}
```

## API参考

### 核心函数

| 函数 | 描述 |
|------|------|
| `Get(names ...string)` | 获取带回退名称的环境变量 |
| `MustGet(names ...string)` | 获取环境变量或未设置时panic |
| `GetOr(name, defaultVal string)` | 获取带显式默认值 |
| `Set(key, value string)` | 设置环境变量 |
| `MustSet(key, value string)` | 设置环境变量或出错时panic |
| `Delete(key string)` | 删除环境变量 |
| `MustDelete(key string)` | 删除环境变量或出错时panic |

### 类型安全获取器

| 函数 | 描述 |
|------|------|
| `GetBool(names ...string)` | 获取布尔值 |
| `GetInt(names ...string)` | 获取整数值 |
| `GetFloat(names ...string)` | 获取浮点值 |
| `Lookup(key string)` | 查找变量存在性 |

### 文件操作

| 函数 | 描述 |
|------|------|
| `LoadFiles(files ...string)` | 加载.env文件 |
| `Expand(value string)` | 在字符串中扩展环境变量 |

### 实用函数

| 函数 | 描述 |
|------|------|
| `Map()` | 将所有环境变量获取为映射 |
| `Key(key string)` | 规范化环境变量键 |
| `Normalize(key string)` | 规范化键带有效性检查 |

## 最佳实践

1. **使用描述性名称**: 选择清晰、描述性的环境变量名称
2. **提供合理默认值**: 始终考虑合理的默认值
3. **早期验证**: 在启动时验证环境变量
4. **记录变量**: 记录应用程序使用的所有环境变量
5. **分组相关变量**: 对相关变量使用一致的命名前缀
6. **避免敏感日志**: 注意不要记录敏感环境变量

## 集成模式

### 与配置

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/config"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

// 环境变量在配置文件中自动替换
cfg := config.Load[Config]()

// 如需要则覆盖环境变量
if host := env.Get("SERVER_HOST"); host != "" {
    cfg.T.Server.Host = host
}
```

### 与功能标志

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/features"
)

// 创建带环境变量默认值的功能标志
var debugMode = features.Bool("debug", env.GetBool("DEBUG"), "启用调试模式")

// 创建带环境变量查找的功能标志
var logLevel = features.String("log_level", env.GetOr("LOG_LEVEL", "info"), "日志级别")
```

### 与日志记录

```go
import (
    "github.com/pubgo/funk/v2/env"
    "github.com/pubgo/funk/v2/log"
)

logger := log.GetLogger("env")

// 记录环境变量加载
env.LoadFiles(".env").Log(logger.RecordErr())

// 记录特定变量访问
host := env.Get("SERVER_HOST")
logger.Info().Str("host", host).Msg("服务器主机已配置")
```