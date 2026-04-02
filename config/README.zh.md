# 配置模块

配置模块提供了一个灵活的系统来管理应用程序配置，支持YAML文件、环境变量替换、CEL表达式引擎和配置合并。

## 功能特性

- **基于YAML的配置**: 使用YAML文件的声明式配置
- **环境变量替换**: 使用`${ENV:-"default_value"}`语法自动替换占位符
- **CEL表达式引擎**: 使用安全沙箱化的CEL表达式动态配置值（`${{env("KEY")}}`）
- **环境变量声明式定义**: 所有环境变量必须在`patch_envs`中预先定义
- **配置合并**: 支持基础和覆盖配置
- **路径安全**: 防止路径遍历攻击
- **线程安全**: 全局配置管理器使用互斥锁保护

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

配置文件支持使用`${ENV:-"default_value"}`语法的环境变量替换：

```yaml
# config.yaml
server:
  host: ${SERVER_HOST:-localhost}
  port: ${SERVER_PORT:-8080}
database:
  url: postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}
```

### CEL表达式引擎

使用安全沙箱化的CEL（Common Expression Language）表达式计算高级配置值。所有表达式使用`${{expression}}`语法：

```yaml
# config.yaml
app:
  config_dir: ${{config_dir()}}
  cert_data: ${{embed("cert.pem")}}
  db_host: ${{env("DB_HOST")}}
  all_envs: ${{envs()}}
```

#### 内置CEL函数

| 函数            | 描述                                         | 示例                       |
| --------------- | -------------------------------------------- | -------------------------- |
| `env("KEY")`    | 获取环境变量值（**必须在patch_envs中定义**） | `${{env("DB_HOST")}}`      |
| `envs()`        | 返回所有已定义环境变量的map                  | `${{envs()}}`              |
| `config_dir()`  | 获取配置文件所在目录                         | `${{config_dir()}}`        |
| `embed("file")` | 将文件内容嵌入为base64（相对于配置目录）     | `${{embed("secret.key")}}` |

#### 环境变量必须预定义

**重要**: 在配置文件中使用`env("KEY")`函数时，该环境变量必须在`patch_envs`中预先定义。这是为了：

1. **显式声明**: 所有配置依赖的环境变量都被明确记录
2. **验证支持**: 使用 `go-playground/validator` 进行灵活的值验证
3. **默认值**: 可以为未设置的环境变量提供默认值
4. **安全性**: 防止意外访问敏感环境变量

示例环境变量定义文件：

```yaml
# configs/envs/database.yaml
DB_HOST:
  desc: 数据库主机地址
  default: localhost

DB_PORT:
  desc: 数据库端口
  default: "5432"
  validate: required,numeric

ADMIN_EMAIL:
  desc: 管理员邮箱
  validate: required,email

API_URL:
  desc: API地址
  validate: url
```

主配置文件引用：

```yaml
# config.yaml
patch_envs:
  - configs/envs/

database:
  host: ${{env("DB_HOST")}}
  port: ${{env("DB_PORT")}}
```

如果在配置中使用了未定义的环境变量，将会得到错误：

```
env: variable "UNDEFINED_VAR" is not defined in patch_envs, all env vars must be declared
```

### 配置结构

```yaml
# config.yaml
resources:
  - configs/database.yaml
  - configs/cache.yaml
patch_resources:
  - configs/local-overrides.yaml
patch_envs:
  - configs/envs/
```

```go
type Resources struct {
    Resources       []string `yaml:"resources"`
    PatchResources  []string `yaml:"patch_resources"`
    PatchEnvs       []string `yaml:"patch_envs"`
}
```

## 核心概念

### 配置加载流程

1. 解析主配置文件的`patch_envs`字段
2. 加载所有环境变量定义（`EnvSpecMap`）
3. 初始化环境变量（验证类型、应用默认值）
4. 解析主配置文件，使用`EnvSpecMap`验证`env()`调用
5. 合并`resources`和`patch_resources`中的配置
6. 返回类型化配置结构

### 环境变量规范（EnvSpec）

```go
type EnvSpec struct {
    Name    string `yaml:"name"`     // 环境变量名称（自动从 key 获取）
    Desc    string `yaml:"desc"`     // 描述
    Default string `yaml:"default"`  // 默认值
    Value   string `yaml:"value"`    // 固定值（优先级高于 Default）
    Example string `yaml:"example"`  // 示例值（仅文档用途）
    Rule    string `yaml:"validate"` // go-playground/validator 验证规则
}
```

#### 常用验证规则

使用 [go-playground/validator](https://github.com/go-playground/validator) 提供的验证规则：

| 规则          | 描述     | 示例                               |
| ------------- | -------- | ---------------------------------- |
| `required`    | 必填     | `validate: required`               |
| `email`       | 邮箱格式 | `validate: email`                  |
| `url`         | URL格式  | `validate: url`                    |
| `uuid`        | UUID格式 | `validate: uuid`                   |
| `ip`          | IP地址   | `validate: ip`                     |
| `ipv4`        | IPv4地址 | `validate: ipv4`                   |
| `ipv6`        | IPv6地址 | `validate: ipv6`                   |
| `numeric`     | 纯数字   | `validate: numeric`                |
| `boolean`     | 布尔值   | `validate: boolean`                |
| `min=n`       | 最小长度 | `validate: min=3`                  |
| `max=n`       | 最大长度 | `validate: max=50`                 |
| `len=n`       | 固定长度 | `validate: len=10`                 |
| `oneof=a b c` | 枚举值   | `validate: oneof=dev staging prod` |

规则可以组合使用，用逗号分隔：

```yaml
APP_NAME:
  desc: 应用名称
  validate: required,min=3,max=50

ENV:
  desc: 运行环境
  default: dev
  validate: oneof=dev staging prod

PORT:
  desc: 服务端口
  default: "8080"
  validate: required,numeric
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

`patch_resources`中的配置会覆盖`resources`中的同名配置。

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
cfg, err := config.LoadFromPath[AppConfig]("path/to/config.yaml")
if err != nil {
  panic(err)
}

// 类型化配置
_ = cfg.T

// 导出最终合并后的 YAML，便于排查与查看
merged, err := config.LoadMergedConfigData("path/to/config.yaml")
if err != nil {
  panic(err)
}
fmt.Println(string(merged))
```

### 注册自定义CEL函数

```go
// 注册自定义函数（必须在加载配置前调用）
err := config.RegisterExpr("custom_func", func(s string) string {
    return strings.ToUpper(s)
})
if err != nil {
    log.Fatal(err)
}

// 在配置中使用
// app:
//   name: ${{custom_func("hello")}}
```

**支持的函数签名**:

| 签名                 | 描述                   |
| -------------------- | ---------------------- |
| `func() T`           | 无参数，返回值         |
| `func() (T, error)`  | 无参数，返回值和错误   |
| `func() error`       | 无参数，只返回错误     |
| `func(T) R`          | 一个参数，返回值       |
| `func(T) (R, error)` | 一个参数，返回值和错误 |
| `func(T) error`      | 一个参数，只返回错误   |

### 路径安全

`embed()`函数会验证路径，防止路径遍历攻击：

```yaml
# ✅ 正常路径
cert: ${{embed("certs/server.pem")}}

# ❌ 路径遍历会被阻止，返回空字符串
secret: ${{embed("../../../etc/passwd")}}
```

## API参考

### 核心函数

| 函数                                                      | 描述                                                |
| --------------------------------------------------------- | --------------------------------------------------- |
| `Load[T]()`                                               | 泛型加载和解析配置                                  |
| `TryLoad[T]()`                                            | 尝试加载类型化配置（返回错误）                      |
| `LoadFromPath[T](cfgPath string)`                         | 从特定路径加载类型化配置                            |
| `LoadMergedConfigData(cfgPath string)`                    | 获取最终合并处理后的配置内容（YAML字节）            |
| `TryLoadMergedData()`                                     | 使用全局配置路径/自动发现尝试加载最终合并配置       |
| `LoadMergedData()`                                        | 加载最终合并配置（错误时 panic）                    |
| `ValidateEnvReferences(cfgPath string)`                   | 手动校验配置中的 env 引用是否在 `patch_envs` 中声明 |
| `SetConfigPath(path string)`                              | 设置自定义配置文件路径                              |
| `GetConfigPath()`                                         | 获取当前配置文件路径                                |
| `GetConfigDir()`                                          | 获取配置目录                                        |
| `GetConfigData(cfgPath string, envSpecMap ...EnvSpecMap)` | 获取处理后的配置数据                                |
| `RegisterExpr(name string, fn any)`                       | 注册自定义CEL表达式函数                             |
| `LoadEnvMap(cfgPath string)`                              | 加载环境变量配置映射                                |

### 配置结构

| 结构         | 描述                                   |
| ------------ | -------------------------------------- |
| `Cfg[T]`     | 包含加载配置和元数据的包装器           |
| `Resources`  | 资源加载和合并的配置                   |
| `EnvSpec`    | 环境变量规范定义                       |
| `EnvSpecMap` | 环境变量规范映射 (map[string]*EnvSpec) |

## 最佳实践

1. **预定义所有环境变量**: 在`patch_envs`中声明所有需要使用的环境变量
2. **使用验证规则**: 使用`validate`字段配置 go-playground/validator 规则
3. **提供默认值**: 为非必需的环境变量指定合理的默认值
4. **使用结构体标签**: 使用结构体标签明确定义YAML映射
5. **早期验证**: 实现配置验证以尽早捕获错误
6. **模块化配置**: 将大配置拆分为逻辑模块
7. **文档化**: 在`patch_envs`中使用`desc`字段记录每个环境变量的用途

## 安全特性

1. **路径遍历防护**: `embed()`函数会验证路径，防止访问配置目录外的文件
2. **沙箱化表达式**: CEL表达式引擎是安全沙箱化的，不允许执行任意代码
3. **敏感数据保护**: 日志中不会输出配置原始内容，防止敏感数据泄露
4. **环境变量声明**: 必须预先定义环境变量，防止意外访问
5. **线程安全**: 全局配置管理器使用读写锁保护

## 从旧版本迁移

如果你从使用`expr-lang/expr`的旧版本迁移，需要更新配置文件语法：

### 旧语法（已弃用）

```yaml
# 旧语法 - 不再支持
app:
  host: ${{env.DB_HOST}}
  debug: ${{env.DEBUG == "true"}}
```

### 新语法（CEL）

```yaml
# 新语法 - CEL表达式
app:
  host: ${{env("DB_HOST")}}
  debug: ${{env("DEBUG") == "true"}}
```

主要变化：
- `env.KEY` → `env("KEY")`
- 所有使用的环境变量必须在`patch_envs`中定义

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