package versioncmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/running"
	"github.com/pubgo/funk/version"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: fmt.Sprintf("%s version info", version.Project()),
		Action: func(ctx context.Context, command *cli.Command) error {
			defer recovery.Exit()
			fmt.Println("project:", version.Project())
			fmt.Println("version:", version.Version())
			fmt.Println("commit-id:", version.CommitID())
			fmt.Println("build-time:", version.BuildTime())
			fmt.Println("instance-id:", running.InstanceID)
			fmt.Println("device-id:", running.DeviceID)
			pretty.Println(running.GetSysInfo())
			return nil
		},
	}
}
