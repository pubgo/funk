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
		name := "feature." + flag.Name
		envVar := env.Key(name)
		options = append(options, redant.Option{
			Flag:        name,
			Description: flag.Usage,
			Value:       flag.Value,
			Default:     flag.Value.String(),
			Category:    category,
			Envs:        []string{envVar},
		})
	})

	return options
}
