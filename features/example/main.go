// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/features"
)

type Data struct {
	Name string
}

func main() {
	// 创建带 Tags 的元数据
	status := features.String("app_status", "ok", "Current application status",
		map[string]any{
			"group":   "health",
			"mutable": true,
		})

	replicas := features.Int("replicas", 1, "Number of replicas",
		map[string]any{
			"group": "scaling",
			"min":   1,
			"max":   10,
		})

	_ = features.String("api_key", "sk-xxxx", "API authentication key",
		map[string]any{
			"sensitive": true,
			"group":     "security",
		})

	_ = features.Json("json_key", &Data{}, "API authentication key",
		map[string]any{
			"sensitive": true,
			"group":     "security",
		})

	// 模拟更新
	go func() {
		time.Sleep(2 * time.Second)
		lo.Must0(status.Set("degraded"))

		time.Sleep(2 * time.Second)
		lo.Must0(replicas.Set(fmt.Sprintf("%v", replicas.Value()+2)))

		time.Sleep(2 * time.Second)
	}()

	// HTTP handler：只输出非敏感字段
	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		data := extractPublicMetadata()
		w.Header().Set("Content-Type", "application/json")
		lo.Must0(json.NewEncoder(w).Encode(data))
	})

	log.Println("Feature server listening on :8181")
	log.Fatal(http.ListenAndServe(":8181", nil))
}

// extractPublicMetadata 导出所有非敏感元数据
func extractPublicMetadata() map[string]any {
	m := make(map[string]any)
	features.VisitAll(func(e *features.Flag) {
		// 跳过敏感字段
		if sensitive, ok := e.Tags["sensitive"].(bool); ok && sensitive {
			m[e.Name] = "******"
			return
		}
		m[e.Name] = e
	})
	return m
}
