package config

import (
	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/vars"
)

func init() {
	vars.Register("config", func() interface{} {
		return typex.M{
			"cfg_type": FileType,
			"cfg_name": FileName,
			"home":     CfgDir,
			"cfg_path": CfgPath,
		}
	})
}
