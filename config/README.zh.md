# 配置模块

配置模块提供了一个灵活的系统来管理应用程序配置，支持YAML文件、环境变量替换和配置合并。

## 功能特性

- **基于YAML的配置**: 使用YAML文件的声明式配置
- **环境变量替换**: 使用`${ENV:"default_value"}`语法自动替换占位符
- **表达式引擎**: 使用类似GitHub Actions工作流语法的表达式动态配置值（`$env.ENV`）
- **配置合并**: 支持基础和覆盖配置
- **热重载**: 无需重启应用程序的动态配置更新

## 安装

```bash
go get github.com/pubgo/funk/v2/config
```

## 快速开始

### 基本配置

```go
import "github.com/pubgo/funk/v2/config"

type ServerConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

type AppConfig struct {
    Server ServerConfig `yaml:"server"`
    Debug  bool         `yaml:"debug"`
}

// 加载配置
cfg := config.Load[AppConfig]()

// 访问配置值
fmt.Printf("服务器: %s:%d\n", cfg.T.Server.Host, cfg.T.Server.Port)
fmt.Printf("调试模式: %t\n", cfg.T.Debug)
```

### 环境变量替换

配置文件支持使用`${ENV:"default_value"}`语法的环境变量替换：

```yaml
# config.yaml
server:
  host: ${SERVER_HOST:"localhost"}
  port: ${SERVER_PORT:8080}
database:
  url: postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}
```

### 表达式引擎

使用类似GitHub Actions工作流语法的表达式计算高级配置值：

```yaml
# config.yaml
app:
  version: ${env.VERSION}
  build_time: ${env.BUILD_TIME}
  config_dir: ${config_dir()}
  cert_data: ${embed("cert.pem")}
```

可用表达式函数：
- `env.ENV_VAR`: 访问环境变量（类似GitHub Actions语法）
- `config_dir()`: 获取配置目录
- `embed(filename)`: 将文件内容嵌入为base64

### 配置结构

```go
type Resources struct {
    Resources       []string `yaml:"resources"`
    PatchResources  []string `yaml:"patch_resources"`
    PatchEnvs       []string `yaml:"patch_envs"`
}
```

## 核心概念

### 配置加载

模块提供泛型`Load[T]()`函数来加载和解析配置文件：

```go
cfg := config.Load[AppConfig]()
```

此函数：
1. 在预定义位置查找配置文件
2. 加载和解析YAML内容
3. 使用`${ENV}`和`${env.ENV}`语法处理环境变量替换
4. 如指定则与附加配置文件合并
5. 返回类型化配置结构

### 环境变量集成

环境变量可在配置文件中使用带可选默认值：

```yaml
# 语法: ${VAR_NAME:"default_value"}
setting1: ${SETTING_1:"default_value"}
setting2: ${REQUIRED_SETTING}  # 无默认值，未设置时出错
```

表达式语法（GitHub Actions风格）：
```yaml
# 语法: ${env.VAR_NAME}
setting1: ${env.SETTING_1}
setting2: ${env.REQUIRED_SETTING}
```

### 配置合并

模块支持合并多个配置文件：

```yaml
# config.yaml
resources:
  - configs/database.yaml
  - configs/cache.yaml
patch_resources:
  - configs/local-overrides.yaml
```

### 表达式计算

可使用表达式计算高级配置值：

```yaml
# config.yaml
app:
  version: ${env.VERSION}
  build_time: ${env.BUILD_TIME}
  config_dir: ${config_dir()}
  cert_data: ${embed("cert.pem")}
```

## 高级用法

### 自定义配置路径

```go
// 设置自定义配置路径
config.SetConfigPath("/path/to/custom/config.yaml")

// 加载配置
cfg := config.Load[AppConfig]()
```

### 手动配置加载

```go
var appConfig AppConfig
envCfgMap := config.LoadFromPath(&appConfig, "path/to/config.yaml")
```

### 配置验证

```go
type ValidatedConfig struct {
    Port int `yaml:"port" validate:"min=1,max=65535"`
}

cfg := config.Load[ValidatedConfig]()
// 添加自定义验证逻辑
```

## API参考

### 核心函数

| 函数 | 描述 |
|------|------|
| `Load[T]()` | 泛型加载和解析配置 |
| `LoadFromPath[T](val *T, cfgPath string)` | 从特定路径加载配置 |
| `SetConfigPath(path string)` | 设置自定义配置文件路径 |
| `GetConfigPath()` | 获取当前配置文件路径 |

### 配置结构

| 结构 | 描述 |
|------|------|
| `Cfg[T]` | 包含加载配置和元数据的包装器 |
| `Resources` | 资源加载和合并的配置 |
| `Node` | 用于灵活访问的YAML节点包装器 |

### 环境集成

| 函数 | 描述 |
|------|------|
| `LoadEnvMap(cfgPath string)` | 加载环境配置映射 |
| `RegisterExpr(name string, expr any)` | 注册自定义表达式函数 |

## 最佳实践

1. **使用结构体标签**: 使用结构体标签明确定义YAML映射
2. **提供默认值**: 为配置选项指定合理的默认值
3. **早期验证**: 实现配置验证以尽早捕获错误
4. **环境覆盖**: 使用环境变量进行部署特定设置
5. **模块化配置**: 将大配置拆分为逻辑模块
6. **文档化**: 记录配置选项及其用途

## 集成模式

### 与功能标志

```go
import (
    "github.com/pubgo/funk/v2/config"
    "github.com/pubgo/funk/v2/features"
)

type AppConfig struct {
    Debug bool `yaml:"debug"`
}

cfg := config.Load[AppConfig]()
debugMode := features.Bool("debug", cfg.T.Debug, "启用调试模式")
```

### 与日志记录

```go
import (
    "github.com/pubgo/funk/v2/config"
    "github.com/pubgo/funk/v2/log"
)

type LogConfig struct {
    Level string `yaml:"level"`
}

cfg := config.Load[LogConfig]()
logger := log.GetLogger("app")
logger.Info().Str("config_level", cfg.T.Level).Msg("配置已加载")
```