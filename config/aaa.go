package config

const (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
	defaultConfigPath = "./configs"
)

var (
	configDir  string
	configPath string
)

type NamedConfig interface {
	// ConfigUniqueName unique name
	ConfigUniqueName() string
}

type Resources struct {
	// Resources resource config file must exist
	Resources []string `yaml:"resources"`

	// PatchResources resource config not required to exist
	PatchResources []string `yaml:"patch_resources"`
}
