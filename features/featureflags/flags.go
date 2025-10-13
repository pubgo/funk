package featureflags

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/features"
	"github.com/urfave/cli/v3"
)

func GetFlags() []cli.Flag {
	var flags []cli.Flag
	features.VisitAll(func(flag *features.Flag) {
		envVar := cli.EnvVars(env.Key("feature." + flag.Name))
		const category = "feature"
		switch flag.Value.Type() {
		case features.BoolType:
			flags = append(flags, &cli.BoolFlag{
				Name:     flag.Name,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(bool),
				Category: category,
				Sources:  envVar,
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, b bool) error {
					return flag.Value.Set(fmt.Sprintf("%v", b))
				},
			})
		case features.StringType:
			flags = append(flags, &cli.StringFlag{
				Name:     flag.Name,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(string),
				Category: category,
				Sources:  envVar,
				Local:    true,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return flag.Value.Set(s)
				},
			})

		case features.JsonType:
			flags = append(flags, &cli.StringFlag{
				Category: category,
				Usage:    flag.Usage,
				Value:    flag.Value.String(),
				Local:    true,
				Sources:  envVar,
				Action: func(ctx context.Context, command *cli.Command, s string) error {
					return flag.Value.Set(s)
				},
			})
		case features.IntType:
			flags = append(flags, &cli.IntFlag{
				Category: category,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(int),
				Local:    true,
				Sources:  envVar,
				Action: func(ctx context.Context, command *cli.Command, i int) error {
					return flag.Value.Set(fmt.Sprintf("%d", i))
				},
			})
		case features.FloatType:
			flags = append(flags, &cli.Float64Flag{
				Category: category,
				Usage:    flag.Usage,
				Value:    flag.Value.Get().(float64),
				Local:    true,
				Sources:  envVar,
				Action: func(ctx context.Context, command *cli.Command, f float64) error {
					return flag.Value.Set(fmt.Sprintf("%f", f))
				},
			})
		}
	})

	return flags
}
