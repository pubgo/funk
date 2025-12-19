package featureflags

import (
	"github.com/pubgo/redant"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/features"
)

func GetFlags() redant.OptionSet {
	const category = "feature"
	var options []redant.Option
	features.VisitAll(func(flag *features.Flag) {
		envVar := env.Key("feature." + flag.Name)
		options = append(options, redant.Option{
			Flag:        flag.Name,
			Description: flag.Usage,
			Value:       flag.Value,
			Default:     flag.Value.String(),
			Category:    category,
			Envs:        []string{envVar},
		})
	})

	return options
}
