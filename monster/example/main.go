// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pubgo/funk/monster"
)

func main() {
	// 创建带 Tags 的元数据
	status := monster.String("app_status", "ok", "Current application status",
		map[string]any{
			"group":   "health",
			"mutable": true,
		})

	replicas := monster.Int("replicas", 1, "Number of replicas",
		map[string]any{
			"group": "scaling",
			"min":   1,
			"max":   10,
		})

	_ = monster.Duration("http_timeout", 5*time.Second, "HTTP request timeout",
		map[string]any{
			"unit":  "seconds",
			"group": "network",
		})

	_ = monster.String("api_key", "sk-xxxx", "API authentication key",
		map[string]any{
			"sensitive": true,
			"group":     "security",
		})

	// 模拟更新
	go func() {
		time.Sleep(2 * time.Second)
		status.Set("degraded")

		time.Sleep(2 * time.Second)
		replicas.Set(replicas.Get() + 2)

		time.Sleep(2 * time.Second)
	}()

	// HTTP handler：只输出非敏感字段
	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		data := extractPublicMetadata()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	log.Println("Monster server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// extractPublicMetadata 导出所有非敏感元数据
func extractPublicMetadata() map[string]any {
	m := make(map[string]any)
	monster.VisitAll(func(e *monster.Entry) {
		// 跳过敏感字段
		if sensitive, ok := e.Tags["sensitive"].(bool); ok && sensitive {
			m[e.Name] = "******"
			return
		}
		m[e.Name] = e.Getter()
	})
	return m
}
