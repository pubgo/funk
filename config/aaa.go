package config

import (
	"github.com/spf13/viper"
)

const (
	defaultConfigName   = "config"
	defaultConfigType   = "yaml"
	defaultConfigPath   = "./configs"
	includeConfigName   = "resources"
	componentConfigKey  = "name"
	defaultComponentKey = "default"
)

type DecoderOption = viper.DecoderConfigOption

type Config interface {
	UnmarshalKey(key string, rawVal interface{}, opts ...DecoderOption) error
	Unmarshal(rawVal interface{}, opts ...DecoderOption) error

	// DecodeComponent decode component config to map[string]*struct
	DecodeComponent(name string, cfgMap interface{}) error
	Get(key string) interface{}
	Set(string, interface{})
	GetString(key string) string
	AllKeys() []string
	All() map[string]interface{}
}
