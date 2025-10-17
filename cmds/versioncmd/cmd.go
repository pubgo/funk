package versioncmd

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/running"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: fmt.Sprintf("%s version info", version.Project()),
		Commands: []*cli.Command{
			{
				Name:  "validate",
				Usage: "show version info",
				Action: func(ctx context.Context, command *cli.Command) error {
					defer recovery.Exit()
					running.CheckVersion()
					return nil
				},
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			defer recovery.Exit()
			pretty.Println(running.GetSysInfo())
			return nil
		},
	}
}
