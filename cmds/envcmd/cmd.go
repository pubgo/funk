package envcmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/recovery"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "envs",
		Usage: "show all envs",
		Action: func(ctx context.Context, command *cli.Command) error {
			defer recovery.Exit()

			env.Reload()

			fmt.Println("config path:", config.GetConfigPath())
			envs := config.LoadEnvConfigMap(config.GetConfigPath())
			for name, cfg := range envs {
				envData := env.Get(name)
				if envData != "" {
					cfg.Default = envData
				}
			}

			pretty.Println(envs)
			return nil
		},
	}
}
