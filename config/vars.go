package config

import (
	"github.com/pubgo/funk/vars"
)

func init() {
	vars.Register("config", func() interface{} {
		return map[string]any{
			"cfg_type": FileType,
			"cfg_name": FileName,
			"home":     CfgDir,
			"cfg_path": CfgPath,
		}
	})
}
