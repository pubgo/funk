// Minimal push-mode profiling example.
//
// Run:
//
//	PYROSCOPE_SERVER=http://localhost:4040 go run ./component/pyroscope/example
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pubgo/funk/v2/component/pyroscope"
)

func main() {
	server := os.Getenv("PYROSCOPE_SERVER")
	if server == "" {
		server = "http://127.0.0.1:4040"
	}

	client := pyroscope.New(pyroscope.Param{
		Cfg: &pyroscope.Config{
			Enabled:         true,
			ServerAddress:   server,
			ApplicationName: "funk.example.pyroscope",
		},
	})
	if !client.Enabled() {
		log.Fatal("pyroscope client disabled")
	}
	defer func() {
		if p := client.Profiler(); p != nil {
			_ = p.Stop()
		}
	}()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		pyroscope.TagWrapper(ctx, pyroscope.Labels("iteration", "demo"), func(ctx context.Context) {
			work()
		})
		time.Sleep(time.Second)
	}

	client.Flush(true)
	log.Println("profiles uploaded")
}

func work() {
	sum := 0
	for i := 0; i < 1_000_000; i++ {
		sum += i
	}
	_ = sum
}
