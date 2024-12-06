package config

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/vars"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
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

func LoadFromPath[T any](val *T, cfgPath string) {
	parentDir := filepath.Dir(cfgPath)
	configBytes := assert.Must1(os.ReadFile(cfgPath))
	configBytes = assert.Must1(envsubst.Bytes(configBytes))

	assert.Must(yaml.Unmarshal(configBytes, val))

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res))

	sort.Strings(res.Resources)
	sort.Strings(res.PatchResources)

	var getRealPath = func([]string) []string {
		var resPaths []string
		for _, resPath := range res.Resources {
			pathList := listAllPath(resPath).Expect("failed to list cfgPath: %s", resPath)
			resPaths = append(resPaths, pathList...)
		}
		resPaths = lo.Filter(resPaths, func(item string, index int) bool { return strings.HasSuffix(item, "."+defaultConfigType) })
		resPaths = lo.Map(resPaths, func(item string, index int) string { return filepath.Join(parentDir, item) })
		return resPaths
	}

	var getCfg = func(resPath string) T {
		resBytes := assert.Must1(os.ReadFile(resPath))
		resBytes = assert.Must1(envsubst.Bytes(resBytes))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1))
		return cfg1
	}

	var cfgList []T
	for _, resPath := range getRealPath(res.Resources) {
		if pathutil.IsNotExist(resPath) {
			log.Panicln("resources config cfgPath not found:", resPath)
		}

		cfgList = append(cfgList, getCfg(resPath))
	}

	for _, resPath := range getRealPath(res.PatchResources) {
		if pathutil.IsNotExist(resPath) {
			continue
		}

		cfgList = append(cfgList, getCfg(resPath))
	}

	assert.Must(Merge(val, cfgList...))
}

func Load[T any]() T {
	if configPath != "" {
		configDir = filepath.Dir(configPath)
	} else {
		configPath, configDir = getConfigPath(defaultConfigName, defaultConfigType)
	}

	var cfg T
	LoadFromPath(&cfg, configPath)
	return cfg
}
