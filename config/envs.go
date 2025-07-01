package config

import (
	"strings"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/samber/lo"
)

type EnvConfigMap map[string]*EnvConf

type EnvConf struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
}

func initEnv(envMap EnvConfigMap) {
	for name, cfg := range envMap {
		envData := env.Get(name)
		envData = strings.TrimSpace(lo.Ternary(envData != "", envData, cfg.Default))
		if cfg.Required && envData == "" {
			panic("env " + cfg.Name + " is required")
		}

		assert.Must(env.Set(name, envData))
	}
}
