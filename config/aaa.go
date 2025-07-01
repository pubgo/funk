package config

type NamedConfig interface {
	// ConfigUniqueName unique name
	ConfigUniqueName() string
}

type Resources struct {
	// Resources resource config file must exist
	Resources []string `yaml:"resources"`

	// PatchResources resource config not required to exist
	PatchResources []string `yaml:"patch_resources"`

	// PatchEnvs envs config not required to exist
	PatchEnvs []string `yaml:"patch_envs"`
}
