package featureflags

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/features"
	"github.com/urfave/cli/v3"
)

func GetFlags() []cli.Flag {
	var flags []cli.Flag
	features.VisitAll(func(flag *features.Flag) {
		switch flag.Value.Type() {
		case "bool":
			flags = append(flags, &cli.BoolFlag{
				Name:     flag.Name,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(bool),
				Category: "feature",
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, b bool) error {
					return flag.Value.Set(fmt.Sprintf("%v", b))
				},
			})
		case "string":
			flags = append(flags, &cli.StringFlag{
				Name:     flag.Name,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(string),
				Category: "feature",
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return flag.Value.Set(s)
				},
			})
		case "json":
			flags = append(flags, &cli.StringFlag{
				Category: "feature",
				Usage:    flag.Usage,
				Value:    flag.Value.String(),
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return flag.Value.Set(s)
				},
			})
		case "int":
			flags = append(flags, &cli.IntFlag{
				Category: "feature",
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(int),
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, i int) error {
					return flag.Value.Set(fmt.Sprintf("%d", i))
				},
			})
		case "float":
			flags = append(flags, &cli.Float64Flag{
				Category: "feature",
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(float64),
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, f float64) error {
					return flag.Value.Set(fmt.Sprintf("%f", f))
				},
			})
		}
	})

	return flags
}
