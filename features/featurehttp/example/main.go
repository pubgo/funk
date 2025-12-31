package main

import (
	"log"

	_ "github.com/pubgo/funk/v2/debugs"
	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/features/featurehttp"
	_ "github.com/pubgo/funk/v2/stack"
)

// Example 展示如何使用 featurehttp 模块
func main() {
	// 创建一些示例 features
	_ = features.String("app_status", "ok", "应用程序状态",
		map[string]any{
			"group":   "health",
			"mutable": true,
		})

	_ = features.Int("replicas", 1, "副本数量",
		map[string]any{
			"group": "scaling",
			"min":   1,
			"max":   10,
		})

	_ = features.Bool("debug", false, "启用调试模式",
		map[string]any{
			"group":   "system",
			"mutable": true,
		})

	_ = features.String("api_key", "sk-xxxx", "API 密钥",
		map[string]any{
			"sensitive": true,
			"group":     "security",
		})

	// 创建并启动 HTTP 服务器
	server := featurehttp.NewServer(":8181")

	// 可选：设置 URL 前缀
	// server.WithPrefix("/features")

	log.Println("Feature HTTP server starting on :8181")
	log.Println("访问 http://localhost:8181/features 查看功能标志管理界面")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// ExampleWithExistingServer 展示如何将 featurehttp 集成到现有的 HTTP 服务器中
func ExampleWithExistingServer() {
	// 创建一些 features
	_ = features.Bool("feature_enabled", true, "功能开关")

	// 创建服务器实例（不启动）
	server := featurehttp.NewServer("")

	// 获取 handler，可以挂载到现有的 mux 上
	handler := server.Handler()

	// 例如，可以这样使用：
	// mux := http.NewServeMux()
	// mux.Handle("/features/", http.StripPrefix("/features", handler))
	// http.ListenAndServe(":8080", mux)

	_ = handler
}
