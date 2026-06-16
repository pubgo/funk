// Runnable CloudEvent demo. Requires a local NATS server with JetStream enabled.
//
// Start NATS (for example with Docker):
//
//	docker run --rm -p 4222:4222 nats:latest -js
//
// Run:
//
//	go run ./component/cloudevent/example
//
// Override NATS URL:
//
//	NATS_URL=nats://127.0.0.1:4222 go run ./component/cloudevent/example
package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samber/lo"
	yaml "gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/component/cloudevent"
	"github.com/pubgo/funk/v2/component/cloudevent/example/demopb"
	"github.com/pubgo/funk/v2/component/lifecycle"
	"github.com/pubgo/funk/v2/component/natsclient"
	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
)

//go:embed config.yaml
var configYAML []byte

type noopLifecycle struct {
	beforeStops []lifecycle.ExecFunc
}

func (l *noopLifecycle) AfterStop(lifecycle.ExecFunc)    {}
func (l *noopLifecycle) BeforeStop(f lifecycle.ExecFunc) { l.beforeStops = append(l.beforeStops, f) }
func (l *noopLifecycle) AfterStart(lifecycle.ExecFunc)   {}
func (l *noopLifecycle) BeforeStart(lifecycle.ExecFunc)  {}

func (l *noopLifecycle) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for i := len(l.beforeStops) - 1; i >= 0; i-- {
		_ = l.beforeStops[i](ctx)
	}
}

func loadConfig() (*cloudevent.Config, error) {
	cfg := new(cloudevent.Config)
	if err := yaml.Unmarshal(configYAML, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}

	lc := new(noopLifecycle)
	nc := natsclient.New(natsclient.Param{
		Cfg: &natsclient.Config{Url: natsURL},
		Lc:  lc,
	})
	jobCli := cloudevent.New(cloudevent.Params{
		Nc:  nc,
		Cfg: cfg,
		Lc:  lc,
	})

	done := make(chan string, 1)
	demopb.RegisterDemoInnerServiceCloudEvent(jobCli, demopb.DemoInnerServiceCloudEvent{
		OnHelloExec: func(ctx context.Context, req *demopb.HelloExecReq) error {
			evt := cloudevent.GetContext(ctx)
			msg := fmt.Sprintf("hello %s on %s", req.GetName(), evt.Subject)
			done <- msg
			return nil
		},
	})

	if err := jobCli.Start(); err != nil {
		log.Fatal(err)
	}
	defer lc.shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pub := demopb.DemoInnerServiceCloudEventPublisher{
		Client: jobCli,
		Opt: func(po *cloudeventpb.PushEventOptions) {
			po.Sender = lo.ToPtr("cloudevent-example")
		},
	}
	if ack := pub.PushHelloExecEvent(ctx, &demopb.HelloExecReq{Name: "world"}); ack.IsErr() {
		log.Fatal(ack.Err())
	} else if info, ok := ack.TryUnwrap(); ok {
		log.Printf("published seq=%d", info.AckInfo.Sequence)
	}

	select {
	case msg := <-done:
		log.Println(msg)
	case <-ctx.Done():
		log.Fatal("timeout waiting for handler")
	}
}
