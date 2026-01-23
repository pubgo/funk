package main

import (
	"context"
	"os"

	"github.com/pubgo/funk/v2/cmds/upgradecmd"
	"github.com/pubgo/redant"
)

func main() {
	// 创建一个用于升级 "pubgo/fastcommit" 的命令
	upgradeCmd := upgradecmd.New("pubgo", "protobuild")

	// 方式1: 将 upgrade 命令作为子命令添加到根命令中
	rootCmd := &redant.Command{
		Use:   "myapp",
		Short: "My application",
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

	// 方式2: 直接运行 upgrade 命令（如果这是唯一的命令）
	// ctx := context.Background()
	// if err := upgradeCmd.Run(ctx); err != nil {
	// 	os.Exit(1)
	// }
}
