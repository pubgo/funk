package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/moby/term"
	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/cliutils"
	"github.com/pubgo/funk/v2/cmds/versioncmd"
	"github.com/pubgo/funk/v2/ctxutil"
	"github.com/pubgo/funk/v2/version"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:                   "testmain",
		Suggest:                true,
		UseShortOptionHandling: true,
		ShellComplete:          cli.DefaultAppComplete,
		Version:                version.Version(),
		Commands: []*cli.Command{
			versioncmd.New(),
		},
		Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
			if !term.IsTerminal(os.Stdin.Fd()) {
				return ctx, fmt.Errorf("stdin is not a terminal")
			}

			if cliutils.IsHelp() {
				return ctx, cli.ShowAppHelp(command)
			}
			return ctx, nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	assert.Exit(app.Run(ctxutil.Signal(), os.Args))
}
