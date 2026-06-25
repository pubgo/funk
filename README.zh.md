# Funk

Funk是一个全面的Go实用程序库，提供增强的错误处理、结果类型、功能标志和各种辅助函数，以简化Go开发。

## 功能特性

### 🚀 增强的错误处理
- 带堆栈跟踪的丰富错误包装
- 与gRPC状态代码兼容的错误代码
- 详细的错误上下文和元数据
- Panic恢复机制

### 📦 结果类型
- 函数式编程启发的结果和错误类型
- 无需显式nil检查的安全错误处理
- 用于映射、过滤和转换结果的综合API
- 支持Future类型的异步操作

### 🎛️ 功能标志
- 运行时动态配置
- 环境变量集成
- CLI标志生成
- 类型安全的功能访问

### 📝 配置管理
- 基于YAML的配置和环境变量替换
- 支持配置合并和扩展
- 用于动态配置值的表达式引擎（类似GitHub Actions工作流语法）
- 热重载功能

### 🪵 结构化日志
- 基于zerolog的高性能日志记录
- 带自动字段注入的上下文感知日志
- 带堆栈跟踪的错误详细信息捕获
- 支持命名空间的模块化日志记录器

### 🌍 环境管理
- 标准化环境变量访问
- 类型安全的环境变量检索
- .env文件加载和解析
- 环境变量扩展

### 📚 堆栈跟踪分析
- 运行时堆栈跟踪捕获
- 调用者识别和元数据提取
- 性能优化的堆栈操作
- 深度Go运行时集成

### 🔧 实用函数
- 用于切片、映射和比较的通用辅助函数
- 断言实用程序
- 路径和文件实用程序
- 字符串和格式化助手

### 🔀 连接分流（同端口多协议）
- 在同一个端口上同时提供 gRPC / HTTP/1 / HTTP/2（通过连接 sniff）
- 包：`connmux`（见 `connmux/README.md`）

## 安装

```bash
go get github.com/pubgo/funk/v2
```

## 快速开始

### 错误处理
```go
import "github.com/pubgo/funk/v2/errors"

err := errors.New("连接失败", errors.Tags{"component": "database"})
err = errors.Wrap(err, "初始化服务")

fmt.Println(err.Error())              // 连接失败
fmt.Println(errors.FormatChain(err))  // 初始化服务: 连接失败
fmt.Println(errors.CollectUserTags(err))
```

配合日志与 result：

```go
import (
    "github.com/pubgo/funk/v2/log"
    "github.com/pubgo/funk/v2/result"
)

log.Error().Err(err).Msg("操作失败") // 自动附加 error_chain、error_tags、error_id

r := result.ErrOf(err)
_ = r.Message()
_ = r.Tags()
```

### 结果类型
```go
import "github.com/pubgo/funk/v2/result"

// 创建成功结果
res := result.OK(42)

// 创建失败结果
errRes := result.Fail[int](errors.New("计算失败"))

// 安全解包值
if value, ok := res.TryUnwrap(); ok {
    fmt.Printf("获得值: %d\n", value)
}

// 链式操作
result := result.OK(10).
    Map(func(x int) int { return x * 2 }).
    FlatMap(func(x int) result.Result[int] {
        if x > 0 {
            return result.OK(x + 1)
        }
        return result.Fail[int](errors.New("负值"))
    })
```

### 功能标志
```go
import "github.com/pubgo/funk/v2/features"

// 定义功能标志
var debugMode = features.Bool("debug", false, "启用调试模式")

// 使用功能标志
if debugMode.Value() {
    log.Println("调试模式已启用")
}
```

### 配置
```go
import "github.com/pubgo/funk/v2/config"

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

cfg := config.Load[Config]()
fmt.Printf("服务器将在 %s:%d 上运行\n", cfg.T.Server.Host, cfg.T.Server.Port)
```

### 日志记录
```go
import "github.com/pubgo/funk/v2/log"

logger := log.GetLogger("myapp")
logger.Info().Msg("应用程序已启动")
logger.Err(someError).Msg("发生错误")
```

### 环境变量
```go
import "github.com/pubgo/funk/v2/env"

// 获取环境变量带默认值
host := env.GetOr("SERVER_HOST", "localhost")

// 类型安全的环境变量访问
port := env.GetInt("SERVER_PORT", "PORT")
debug := env.GetBool("DEBUG")
```

### 堆栈跟踪分析
```go
import "github.com/pubgo/funk/v2/stack"

// 捕获当前堆栈跟踪
frames := stack.Trace()

// 获取调用者信息
caller := stack.Caller(0)
fmt.Printf("被调用于: %s:%d\n", caller.File, caller.Line)
```

## 模块

- **errors**: 带包装、堆栈跟踪、标签和可读错误链的增强错误处理
- **result**: 用于更安全错误处理的函数式结果和错误类型
- **features**: 用于运行时配置的功能标志系统
- **assert**: 用于测试和验证的断言实用程序
- **config**: 支持多种来源的配置管理
- **log**: 增强的日志记录功能
- **env**: 环境变量管理
- **stack**: 堆栈跟踪分析和调用者识别

## 文档

有关详细文档，请访问：
- [错误处理](./errors/README.md)
- [错误码](./errors/errcode/README.md)
- [结果类型](./result/README.md)
- [功能标志](./features/README.md)
- [配置](./config/README.md)
- [日志记录](./log/README.md)
- [环境](./env/README.md)
- [堆栈](./stack/README.md)

- [发布指南](./docs/RELEASE.md) — 最新版本：[v2.0.5](./docs/releases/v2.0.5.md)

## 贡献

欢迎贡献！在提交拉取请求之前，请阅读我们的贡献指南。

## 许可证

MIT