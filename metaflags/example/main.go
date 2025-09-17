package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pubgo/funk/metaflags"
)

var startTime = time.Now()

func main() {
	status := metaflags.String("app_status", "ok", "Current app status")
	replicas := metaflags.Int("replicas", 1, "Number of replicas")
	timeout := metaflags.Duration("timeout", 5*time.Second, "Request timeout")
	features := metaflags.StringSlice("features", []string{"auth"}, "Enabled features")
	labels := metaflags.StringMap("labels", map[string]string{"env": "dev"}, "Node labels")

	// 模拟更新
	go func() {
		time.Sleep(2 * time.Second)
		status.Set("degraded")

		time.Sleep(2 * time.Second)
		replicas.Set(3)

		time.Sleep(2 * time.Second)
		timeout.Set(10 * time.Second)

		time.Sleep(2 * time.Second)
		features.Set([]string{"auth", "metrics", "tracing"})

		time.Sleep(2 * time.Second)
		lbls := labels.Get().(map[string]string)
		lbls["updated"] = time.Now().Format("15:04")
		labels.Set(lbls)
	}()

	// HTTP 输出所有元数据
	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{")
		first := true
		metaflags.VisitAll(func(e *metaflags.Entry) {
			if !first {
				fmt.Fprint(w, ",")
			}
			first = false
			fmt.Fprintf(w, "\n  %q: %v", e.Name, mustJSON(e.Value.Get()))
		})
		fmt.Fprintln(w, "\n}")
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 简单 JSON 转换（生产环境用 json.Marshal）
func mustJSON(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case []string:
		return fmt.Sprintf("%q", val)
	case map[string]string:
		var parts []string
		for k, v := range val {
			parts = append(parts, fmt.Sprintf("%q:%q", k, v))
		}
		return fmt.Sprintf("{%s}", strings.Join(parts, ","))
	default:
		return fmt.Sprintf("%v", val)
	}
}
