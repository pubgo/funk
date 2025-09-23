package config

import (
	"strings"

	"github.com/pubgo/funk/env"
	"github.com/samber/lo"
)

type EnvSpecMap map[string]*EnvSpec

type EnvSpec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
	Example     string `yaml:"example"`
}

func initEnv(envMap EnvSpecMap) {
	for name, cfg := range envMap {
		envData := env.Get(name)
		envData = strings.TrimSpace(lo.Ternary(envData != "", envData, cfg.Default))
		if cfg.Required && envData == "" {
			panic("env " + cfg.Name + " is required")
		}

		env.Set(name, envData).Must()
	}
}
