package config

import (
	"strings"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/samber/lo"
)

type EnvConfigMap map[string]*EnvConf

type EnvConf struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
	Example     string `yaml:"example"`
	Versions    string `yaml:"versions"`
	Tags        string `yaml:"tags"`
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
