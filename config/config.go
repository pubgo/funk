package config

import (
	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/vars"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
	defaultConfigPath = "./configs"
)

var (
	configDir  string
	configPath string
)

func init() {
	vars.RegisterValue("config", map[string]any{
		"config_type": defaultConfigType,
		"config_name": defaultConfigName,
		"config_path": configPath,
		"config_dir":  configDir,
	})
}

func LoadFromPath[T any](val *T, path string) {
	dir := filepath.Dir(path)
	configBytes := assert.Must1(os.ReadFile(path))
	configBytes = assert.Must1(envsubst.Bytes(configBytes))

	assert.Must(yaml.Unmarshal(configBytes, val))

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res))

	var cfgList []T
	for _, resPath := range res.Resources {
		resAbsPath := filepath.Join(dir, resPath)
		if pathutil.IsNotExist(resAbsPath) {
			log.Panicln("resources config path not found:", resAbsPath)
		}

		resBytes := assert.Must1(os.ReadFile(resAbsPath))
		resBytes = assert.Must1(envsubst.Bytes(resBytes))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1))
		cfgList = append(cfgList, cfg1)
	}

	for _, resPath := range res.PatchResources {
		resAbsPath := filepath.Join(dir, resPath)
		if pathutil.IsNotExist(resAbsPath) {
			continue
		}

		resBytes := assert.Must1(os.ReadFile(resAbsPath))
		resBytes = assert.Must1(envsubst.Bytes(resBytes))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1))
		cfgList = append(cfgList, cfg1)
	}

	assert.Must(Merge(val, cfgList...))
}

func Load[T any]() T {
	if configPath != "" {
		configDir = filepath.Dir(configPath)
	} else {
		configPath, configDir = getConfigPath(defaultConfigName, defaultConfigType)
	}

	configBytes := assert.Must1(os.ReadFile(configPath))
	configBytes = assert.Must1(envsubst.Bytes(configBytes))

	var cfg T
	assert.Must(yaml.Unmarshal(configBytes, &cfg))

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res))

	var cfgList []T
	for _, resPath := range res.Resources {
		resAbsPath := filepath.Join(configDir, resPath)
		if pathutil.IsNotExist(resAbsPath) {
			log.Panicln("resources config path not found:", resAbsPath)
		}

		resBytes := assert.Must1(os.ReadFile(resAbsPath))
		resBytes = assert.Must1(envsubst.Bytes(resBytes))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1))
		cfgList = append(cfgList, cfg1)
	}

	for _, resPath := range res.PatchResources {
		resAbsPath := filepath.Join(configDir, resPath)
		if pathutil.IsNotExist(resAbsPath) {
			continue
		}

		resBytes := assert.Must1(os.ReadFile(resAbsPath))
		resBytes = assert.Must1(envsubst.Bytes(resBytes))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1))
		cfgList = append(cfgList, cfg1)
	}

	assert.Must(Merge(&cfg, cfgList...))
	return cfg
}
