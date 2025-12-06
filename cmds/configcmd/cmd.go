package configcmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	yaml "gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/recovery"
)

func New[Cfg any]() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "config management",
		Commands: []*cli.Command{
			{
				Name:        "show",
				Description: "show config data",
				Action: func(ctx context.Context, command *cli.Command) error {
					defer recovery.Exit()
					fmt.Println("config path:\n", config.GetConfigPath())
					fmt.Println("config raw data:\n", string(assert.Must1(yaml.Marshal(config.Load[Cfg]().T))))
					return nil
				},
			},
		},
	}
}
