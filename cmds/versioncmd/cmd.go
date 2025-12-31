package versioncmd

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/running"
	"github.com/pubgo/redant"
)

func New() *redant.Command {
	return &redant.Command{
		Use:   "version",
		Short: fmt.Sprintf("%s version info", version.Project()),
		Children: []*redant.Command{
			{
				Use:   "validate",
				Short: "show version info",
				Handler: func(ctx context.Context, i *redant.Invocation) error {
					defer recovery.Exit()
					running.CheckVersion()
					return nil
				},
			},
		},
		Handler: func(ctx context.Context, i *redant.Invocation) error {
			defer recovery.Exit()
			pretty.Println(running.GetSysInfo())
			return nil
		},
	}
}
