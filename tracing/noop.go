package tracing

import (
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/recovery"
)

func init() {
	defer recovery.Exit()

	RegisterFactory("noop", func(cfg config.CfgMap) error { return nil })
}
