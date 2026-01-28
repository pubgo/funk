package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/strutil"
)

var (
	validate     *validator.Validate
	validateOnce sync.Once
)

// getValidator returns a singleton validator instance
func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()
	})
	return validate
}

type EnvSpecMap map[string]*EnvSpec

// ValidateAll validates all environment specs in the map
func (m EnvSpecMap) ValidateAll() error {
	var errs []error
	for name, spec := range m {
		spec.Name = name
		if err := spec.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Errorf("env validation failed: %v", errs)
	}
	return nil
}

type EnvSpec struct {
	// Name of the environment variable.
	Name string `yaml:"name"`

	Desc    string `yaml:"desc"`
	Default string `yaml:"default"`
	Value   string `yaml:"value"`
	Example string `yaml:"example"`

	// Validation using go-playground/validator tags
	// Examples: "required", "email", "url", "uuid", "ip", "numeric", "min=3,max=50"
	Rule string `yaml:"validate"`
}

func (e EnvSpec) GetValue() string {
	return strings.TrimSpace(strutil.FirstNotEmpty(env.Get(e.Name), e.Value, e.Default))
}

// Validate performs validation on the env spec using go-playground/validator
func (e EnvSpec) Validate() error {
	value := e.GetValue()

	// Skip validation if no validate rule
	if e.Rule == "" {
		return nil
	}

	// Use go-playground/validator for all validation
	v := getValidator()
	if err := v.Var(value, e.Rule); err != nil {
		return fmt.Errorf("env %q: validation failed for value %q with rule %q: %v", e.Name, value, e.Rule, err)
	}

	return nil
}

// initEnv initializes environment variables from the spec map.
// Returns error if validation fails instead of panicking.
func initEnv(envMap EnvSpecMap) error {
	if envMap == nil {
		return nil
	}

	// First pass: set names and validate all
	var errs []error
	for name, cfg := range envMap {
		cfg.Name = name
		if err := cfg.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Err(err).Msg("env validation failed")
		}
		return errors.Errorf("env validation failed with %d errors", len(errs))
	}

	// Second pass: set env values
	for name, cfg := range envMap {
		env.Set(name, cfg.GetValue()).MustWithLog()
	}

	return nil
}
