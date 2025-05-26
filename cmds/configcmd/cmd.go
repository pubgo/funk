package configcmd

import (
	"context"
	"fmt"

	"github.com/pubgo/dix"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func New[Cfg any](di *dix.Dix) *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "config management",
		Commands: []*cli.Command{
			{
				Name:        "show",
				Description: "show config data",
				Action: func(ctx context.Context, command *cli.Command) error {
					fmt.Println("config path:", config.GetConfigPath())
					fmt.Println("config raw data:", string(assert.Must1(yaml.Marshal(config.Load[Cfg]().T))))
					return nil
				},
			},
		},
	}
}
