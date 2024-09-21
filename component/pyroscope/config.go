package pyroscope

type Config struct {
	Enabled       bool   `yaml:"enabled"`
	ServerAddress string `yaml:"server_address"`
}
