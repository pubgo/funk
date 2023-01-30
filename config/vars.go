package config

import (
	"github.com/pubgo/funk/vars"
)

func init() {
	vars.Register("config", getCfgData)
}
