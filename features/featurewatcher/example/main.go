package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pubgo/funk/v2/closer"
	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/features/featurewatcher"
)

// Example 展示如何使用 featurewatcher
func main() {
	// 创建一些功能标志
	_ = features.String("app_status", "ok", "应用程序状态",
		map[string]any{
			"group":   "health",
			"mutable": true,
		})

	_ = features.Int("replicas", 1, "副本数量",
		map[string]any{
			"group":   "scaling",
			"mutable": true,
		})

	_ = features.Bool("debug", false, "启用调试模式",
		map[string]any{
			"group":   "system",
			"mutable": true,
		})

	// 创建监控目录
	watchDir := ".local/features"
	if err := os.MkdirAll(watchDir, 0755); err != nil {
		log.Fatal(err)
	}

	// 创建示例配置文件
	configFile := filepath.Join(watchDir, "config.txt")
	exampleContent := `# Feature Flags 配置文件
# 格式: feature_name=feature_value
# 注释行以 # 或 // 开头

app_status=ok
replicas=3
debug=true
`

	if err := os.WriteFile(configFile, []byte(exampleContent), 0644); err != nil {
		log.Fatal(err)
	}

	// 创建监控器
	watcher, err := featurewatcher.NewWatcher(watchDir,
		featurewatcher.WithDebounce(500*time.Millisecond),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 启动监控
	if err := watcher.Start(); err != nil {
		log.Fatal(err)
	}
	defer closer.ErrClose(watcher.Stop)

	log.Println("Feature watcher started, monitoring:", watchDir)
	log.Println("You can edit files in", watchDir, "to update features")
	log.Println("Press Ctrl+C to stop")

	// 保持运行
	select {}
}

// ExampleWithCustomFeature 展示如何使用自定义的 Feature 实例
func mainWithCustomFeature() {
	// 创建自定义 Feature 实例
	customFeature := features.NewFeature()

	_ = customFeature.AddFunc("custom_flag", "自定义标志", &customValue{val: "default"}, nil)

	// 创建监控器，使用自定义 Feature
	watcher, err := featurewatcher.NewWatcher("./features",
		featurewatcher.WithFeature(customFeature),
		featurewatcher.WithDebounce(1*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := watcher.Start(); err != nil {
		log.Fatal(err)
	}
	defer closer.ErrClose(watcher.Stop)

	select {}
}

// 示例自定义 Value 实现
type customValue struct {
	val string
}

func (v *customValue) String() string { return v.val }
func (v *customValue) Set(s string) error {
	v.val = s
	return nil
}
func (v *customValue) Type() string { return "string" }
func (v *customValue) Value() any   { return v.val }
