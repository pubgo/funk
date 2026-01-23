# upgradecmd

通用自升级命令组件，用于 CLI 应用程序自动升级功能。

## 功能特性

- 从 GitHub Releases 自动下载最新版本
- 支持跨平台（Windows, macOS, Linux）
- 支持多架构（amd64, arm64 等）
- 交互式版本选择
- 下载进度显示
- 自动替换可执行文件

## 使用方法

### 方式1: 作为子命令添加到根命令

```go
package main

import (
    "context"
    "os"
    
    "github.com/pubgo/funk/v2/cmds/upgradecmd"
    "github.com/pubgo/redant"
)

func main() {
    // 创建升级命令，指定 GitHub owner 和 repo
    upgradeCmd := upgradecmd.New("your-github-username", "your-repo-name")

    // 将命令添加到你的 CLI 应用中
    rootCmd := &redant.Command{
        Use:   "your-app",
        Short: "Your application description",
        Children: []*redant.Command{
            upgradeCmd,
            // 可以添加其他命令...
        },
    }

    // 运行根命令
    ctx := context.Background()
    if err := rootCmd.Run(ctx); err != nil {
        os.Exit(1)
    }
}
```

### 方式2: 直接运行 upgrade 命令

```go
package main

import (
    "context"
    "os"
    
    "github.com/pubgo/funk/v2/cmds/upgradecmd"
)

func main() {
    // 创建升级命令
    upgradeCmd := upgradecmd.New("your-github-username", "your-repo-name")

    // 直接运行命令
    ctx := context.Background()
    if err := upgradeCmd.Run(ctx); err != nil {
        os.Exit(1)
    }
}
```

## 目录结构

```
cmds/upgradecmd/
├── cmd.go           # 主要命令实现
├── progress.go      # 下载进度条实现
├── githubclient/    # GitHub API 客户端
│   ├── asset.go     # Asset 结构和解析
│   ├── release.go   # Release API 封装
│   └── utils.go     # 工具函数
├── examples/        # 使用示例
│   └── example.go   # 基本使用示例
└── README.md        # 本文档
```

## 命令用法

```bash
# 列出所有可用版本
your-app upgrade list

# 交互式选择并升级到指定版本
your-app upgrade
```

## 依赖要求

- 应用程序需要在 GitHub 上有 releases
- Release assets 需要按照约定命名，包含操作系统和架构信息
- 例如：`myapp-v1.0.0-linux-amd64`, `myapp-v1.0.0-darwin-arm64`

## 架构支持

自动检测并支持以下平台：
- Linux (amd64, arm64, 386)
- macOS/Darwin (amd64, arm64)
- Windows (amd64, 386)
