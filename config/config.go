package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
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

func GetConfigData(cfgPath string) (_ []byte, gErr error) {
	var configBytes []byte
	defer recovery.Err(&gErr, func(err error) error {
		log.Err(err).Str("config_path", cfgPath).Msgf("config: %s", configBytes)
		return err
	})

	configBytes = result.Of(os.ReadFile(cfgPath)).Expect("failed to read config data: %s", cfgPath)
	configBytes = result.Of(envsubst.Bytes(configBytes)).Expect("failed to handler config env data: %s", cfgPath)
	configBytes = cfgFormat(configBytes, &config{workDir: filepath.Dir(cfgPath)})
	return configBytes, nil
}

func LoadFromPath[T any](val *T, cfgPath string) EnvConfigMap {
	defer recovery.Exit(func(err error) error {
		log.Err(err).Str("config_path", cfgPath).Msg("failed to load config")
		return err
	})

	valType := reflect.TypeOf(val)
	for {
		if valType.Kind() != reflect.Ptr {
			break
		}

		valType = valType.Elem()
	}
	if valType.Kind() != reflect.Struct {
		log.Panic().
			Str("config_path", cfgPath).
			Str("type", fmt.Sprintf("%#v", val)).
			Msg("config type not correct")
	}

	configBytes := result.Of(GetConfigData(cfgPath)).Expect("failed to handler config data")
	defer recovery.Exit(func(err error) error {
		log.Err(err).
			Str("config_path", cfgPath).
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
		return nil
	}

	parentDir := filepath.Dir(cfgPath)
	getRealPath := func(pp []string) []string {
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
	getCfg := func(resPath string) T {
		resBytes := result.Of(GetConfigData(resPath)).Expect("failed to handler config data")

		var cfg1 T
		result.Err[any](yaml.Unmarshal(resBytes, &cfg1)).
			Unwrap(func(err error) error {
				fmt.Println("res_path", resPath)
				fmt.Println("config_data", string(resBytes))
				assert.Exit(os.WriteFile(resPath+".err.yml", resBytes, 0666))
				return errors.Wrap(err, "failed to unmarshal config")
			})

		return cfg1
	}

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res), "failed to unmarshal resource config")

	var envCfgMap EnvConfigMap
	for _, envPath := range res.PatchEnvs {
		envPath = filepath.Join(parentDir, envPath)
		if pathutil.IsNotExist(envPath) {
			log.Warn().Str("env_path", envPath).Msg("env config cfgPath not found")
			continue
		}

		pathList := listAllPath(envPath).Expect("failed to list envPath: %s", envPath)
		for _, p := range pathList {
			envConfigBytes := result.Of(os.ReadFile(p)).Expect("failed to handler env config data, path=%s", p)
			assert.MustF(yaml.Unmarshal(envConfigBytes, &envCfgMap), "failed to unmarshal env config, path=%s", p)
		}
	}
	initEnv(envCfgMap)

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
	return envCfgMap
}

type Cfg[T any] struct {
	T      T
	P      *T
	EnvCfg *EnvConfigMap
}

func Load[T any]() Cfg[T] {
	if configPath != "" {
		configDir = filepath.Dir(configPath)
	} else {
		configPath, configDir = getConfigPath(defaultConfigName, defaultConfigType)
	}

	var cfg T
	cfgMap := LoadFromPath(&cfg, configPath)
	return Cfg[T]{T: cfg, P: &cfg, EnvCfg: lo.ToPtr(cfgMap)}
}
