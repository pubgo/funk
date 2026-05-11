package configflags

import (
	"github.com/pubgo/redant"
	"github.com/samber/lo"
	"github.com/spf13/pflag"

	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/env"
)

var ConfFlag = redant.Option{
	Flag:        "config",
	Shorthand:   "c",
	Description: "config path",
	Default:     config.GetConfigPath(),
	Value:       redant.StringOf(lo.ToPtr(config.GetConfigPath())),
	Inherit:     true,
	Envs:        []string{env.Key("config_path")},
	Action: func(val pflag.Value) error {
		config.SetConfigPath(val.String())
		env.Set("config_path", val.String())
		return nil
	},
}
