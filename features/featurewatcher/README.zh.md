# Feature Watcher 模块

Feature Watcher 模块提供了文件监控功能，可以监控指定目录中的文件变化，并自动将文件内容解析并应用到功能标志（Feature Flags）中。

## 功能特性

- 📁 **目录监控**: 监控指定目录中的所有文件
- 🔄 **自动重载**: 文件变化时自动重新加载并应用
- 📝 **简单格式**: 支持 `name=value` 格式的配置文件
- ⚡ **防抖处理**: 避免频繁触发，提高性能
- 🔒 **安全控制**: 自动跳过不可修改的功能标志
- 📊 **详细日志**: 记录加载过程和错误信息

## 文件格式

配置文件使用简单的 `name=value` 格式，每行一个功能标志：

```
# 这是注释，以 # 或 // 开头
# 空行会被忽略

app_status=ok
replicas=3
debug=true
log_level=info
max_retries=5
```

## 快速开始

### 基本使用

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/features/featurewatcher"
)

func main() {
	// 注册功能标志
	_ = features.String("app_status", "ok", "应用程序状态")
	_ = features.Int("replicas", 1, "副本数量")
	_ = features.Bool("debug", false, "启用调试模式")

	// 创建监控目录
	watchDir := "./features"
	os.MkdirAll(watchDir, 0755)

	// 创建监控器
	watcher, err := featurewatcher.NewWatcher(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	// 启动监控
	if err := watcher.Start(); err != nil {
		log.Fatal(err)
	}
	defer watcher.Stop()

	log.Println("监控已启动，编辑", watchDir, "目录中的文件来更新功能标志")
	
	// 保持运行
	select {}
}
```

### 自定义配置

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithDebounce(1*time.Second), // 设置防抖时间为 1 秒
)
```

### 使用自定义 Feature 实例

```go
customFeature := features.NewFeature()
_ = customFeature.AddFunc("custom_flag", "自定义标志", value, nil)

watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithFeature(customFeature), // 使用自定义 Feature 实例
)
```

## 配置选项

### WithFeature(f *features.Feature)

指定要使用的 Feature 实例。如果不指定，默认使用全局的 `defaultFeature`。

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithFeature(customFeature),
)
```

### WithDebounce(d time.Duration)

设置防抖时间，避免文件频繁变化时触发多次加载。默认值为 500 毫秒。

```go
watcher, err := featurewatcher.NewWatcher("./features",
	featurewatcher.WithDebounce(1*time.Second),
)
```

## 工作原理

1. **初始加载**: 启动时自动加载目录中的所有文件
2. **文件监控**: 使用 `fsnotify` 监控目录中的文件变化
3. **事件处理**: 当文件被写入或创建时，触发重新加载
4. **防抖处理**: 在防抖时间内的多次变化只处理最后一次
5. **解析应用**: 解析文件内容，查找对应的功能标志并更新值

## 文件解析规则

1. **格式**: 每行必须是 `name=value` 格式
2. **注释**: 以 `#` 或 `//` 开头的行会被忽略
3. **空行**: 空行会被忽略
4. **空格**: 名称和值前后的空格会被自动去除
5. **错误处理**: 格式错误的行会被记录但不会中断处理

## 安全特性

### 自动跳过不可修改的功能标志

如果功能标志标记为 `mutable: false`，即使文件中包含该标志，也不会被更新：

```go
_ = features.String("version", "1.0.0", "版本号",
	map[string]any{
		"mutable": false, // 不可修改
	})
```

### 自动跳过未注册的功能标志

如果文件中包含未注册的功能标志名称，会被记录但不会报错。

## 日志输出

模块会输出详细的日志信息：

- **INFO**: 监控启动/停止、文件加载成功
- **DEBUG**: 功能标志更新详情
- **WARN**: 格式错误、未找到功能标志等警告
- **ERROR**: 文件读取错误、设置值失败等错误

## 注意事项

1. **目录必须存在**: 监控的目录必须存在，否则会返回错误
2. **文件格式**: 文件必须使用 `name=value` 格式，其他格式会被忽略
3. **功能标志必须先注册**: 只有已注册的功能标志才会被更新
4. **并发安全**: 文件加载和功能标志更新都是线程安全的
5. **防抖时间**: 合理设置防抖时间可以避免频繁触发，但也会增加延迟

## 示例场景

### 场景 1: 开发环境动态配置

在开发环境中，可以通过修改配置文件来动态调整功能标志，无需重启应用。

### 场景 2: 多环境配置

不同环境可以使用不同的配置文件，通过文件监控实现配置的动态切换。

### 场景 3: 配置中心集成

可以配合配置中心（如 etcd、Consul）的文件同步功能，实现配置的自动更新。

## 错误处理

模块会处理以下错误情况：

- 目录不存在
- 文件读取失败
- 格式错误（记录警告，继续处理其他行）
- 功能标志不存在（记录调试信息，跳过）
- 设置值失败（记录错误，继续处理其他行）

所有错误都会记录到日志中，不会中断监控过程。

