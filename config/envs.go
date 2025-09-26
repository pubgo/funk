package config

import (
	"strings"

	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/strutil"
)

type EnvSpecMap map[string]*EnvSpec

type EnvSpec struct {
	Name string `yaml:"name"`

	// Description Deprecated: use Desc instead.
	Description string `yaml:"description"`
	Desc        string `yaml:"desc"`
	Default     string `yaml:"default"`
	Value       string `yaml:"value"`
	Required    bool   `yaml:"required"`
	Example     string `yaml:"example"`
}

func initEnv(envMap EnvSpecMap) {
	for name, cfg := range envMap {
		cfg.Name = name

		envData := strings.TrimSpace(strutil.FirstNotEmpty(env.Get(name), cfg.Value, cfg.Default))
		if cfg.Required && envData == "" {
			panic("env " + cfg.Name + " is required")
		}

		env.Set(name, envData).Must()
	}
}
