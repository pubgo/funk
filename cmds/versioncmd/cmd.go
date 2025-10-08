package versioncmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/pubgo/funk/v2/buildinfo"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/running"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: fmt.Sprintf("%s version info", buildinfo.Project()),
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
