package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/typex"
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
	defer recovery.Exit(func(err error) error {
		log.Err(err).Str("config_path", cfgPath).Msg("failed to load config")
		return err
	})

	parentDir := filepath.Dir(cfgPath)
	configBytes := result.Of(os.ReadFile(cfgPath)).Expect("failed to read config data: %s", cfgPath)
	configBytes = result.Of(envsubst.Bytes(configBytes)).Expect("failed to handler config env data: %s", cfgPath)

	defer recovery.Exit(func(err error) error {
		log.Err(err).
			Str("config_path", cfgPath).
			Str("config_dir", parentDir).
			Str("config_data", string(configBytes)).
			Msg("failed to load config")
		return err
	})

	if err := yaml.Unmarshal(configBytes, val); err != nil {
		log.Panic().
			Err(err).
			Str("config_data", string(configBytes)).
			Str("config_path", cfgPath).
			Msg("failed to unmarshal config")
		return
	}

	var getRealPath = func(pp []string) []string {
		pp = lo.Map(pp, func(item string, index int) string { return filepath.Join(parentDir, item) })

		var resPaths []string
		for _, resPath := range pp {
			pathList := listAllPath(resPath).Expect("failed to list cfgPath: %s", resPath)
			resPaths = append(resPaths, pathList...)
		}

		// skip .*.yaml and cfg.other
		var cfgFilter = func(item string, index int) bool {
			return strings.HasSuffix(item, "."+defaultConfigType) && !strings.HasPrefix(item, ".")
		}
		resPaths = lo.Filter(resPaths, cfgFilter)
		return lo.Uniq(resPaths)
	}
	var getCfg = func(resPath string) T {
		resBytes := result.Of(os.ReadFile(resPath)).Expect("failed to read config data: %s", resPath)
		resBytes = result.Of(envsubst.Bytes(resBytes)).Expect("failed to handler config env data: %s", resPath)
		resBytes = []byte(cfgFormat(string(resBytes), &config{
			workDir: filepath.Dir(resPath),
		}))

		var cfg1 T
		assert.Must(yaml.Unmarshal(resBytes, &cfg1), "failed to unmarshal config")
		return cfg1
	}

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res), "failed to unmarshal resource config")

	var cfgList []T
	cfgList = append(cfgList, typex.DoBlock1(func() []T {
		var resPathList = getRealPath(res.Resources)
		sort.Strings(resPathList)

		var pathList []T
		for _, resPath := range resPathList {
			if pathutil.IsNotExist(resPath) {
				log.Panic().Str("path", resPath).Msg("resources config cfgPath not found")
				continue
			}

			pathList = append(pathList, getCfg(resPath))
		}
		return pathList
	})...)
	cfgList = append(cfgList, typex.DoBlock1(func() []T {
		var patchResPathList = getRealPath(res.PatchResources)
		sort.Strings(patchResPathList)

		var pathList []T
		for _, resPath := range patchResPathList {
			if pathutil.IsNotExist(resPath) {
				continue
			}

			pathList = append(pathList, getCfg(resPath))
		}
		return pathList
	})...)

	assert.Exit(Merge(val, cfgList...), "failed to merge config")
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
