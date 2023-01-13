package config

import (
	"github.com/pubgo/funk/env"
)

func init() {
	flags.Register(&cli.StringFlag{
		Name:        "home",
		Destination: &CfgPath,
		Usage:       "config home dir, [configs]",
		EnvVars:     typex.StrOf(env.Key(consts.EnvHome)),
	})

	flags.Register(&cli.StringFlag{
		Name:  "config",
		Usage: "config file name",
		Value: FileName + "." + FileType,
	})
}
