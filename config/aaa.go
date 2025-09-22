package config

type NamedConfig interface {
	// ConfigUniqueName unique name
	ConfigUniqueName() string
}

type Resources struct {
	// Resources resource config file or dir must exist
	Resources []string `yaml:"resources"`

	// PatchResources resource config file or dir not required to exist
	PatchResources []string `yaml:"patch_resources"`

	// PatchEnvs env config file or dir not required to exist
	PatchEnvs []string `yaml:"patch_envs"`
}
