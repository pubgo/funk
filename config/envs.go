package config

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/strutil"
)

type EnvSpecMap map[string]*EnvSpec

type EnvSpec struct {
	// Name of the environment variable.
	Name string `yaml:"name"`

	// Description Deprecated: use Desc instead.
	Description string `yaml:"description"`
	Desc        string `yaml:"desc"`
	Default     string `yaml:"default"`
	Value       string `yaml:"value"`
	Required    bool   `yaml:"required"`
	Example     string `yaml:"example"`
}

func (e EnvSpec) GetValue() string {
	return strings.TrimSpace(strutil.FirstNotEmpty(env.Get(e.Name), e.Value, e.Default))
}

func (e EnvSpec) Validate() error {
	if e.Required && e.GetValue() == "" {
		return fmt.Errorf("env:%s is required", e.Name)
	}
	return nil
}

func initEnv(envMap EnvSpecMap) {
	for name, cfg := range envMap {
		cfg.Name = name

		lo.Must0(cfg.Validate())
		env.Set(name, cfg.GetValue()).MustWithLog()
	}
}
