package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/logfields"
	"github.com/pubgo/funk/v2/pathutil"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/typex"
	"github.com/pubgo/funk/v2/vars"
)

const (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
	defaultConfigPath = "./configs"
)

func init() {
	vars.Register("config", func() any {
		return map[string]any{
			"config_type": defaultConfigType,
			"config_name": defaultConfigName,
			"config_path": globalManager.GetPath(),
			"config_dir":  globalManager.GetDir(),
		}
	})
}

func GetConfigDir() string {
	return globalManager.GetDir()
}

func GetConfigPath() string {
	return globalManager.GetPath()
}

func SetConfigPath(confPath string) {
	assert.If(confPath == "", "config path is null")
	globalManager.SetPath(confPath)
}

// GetConfigData loads and processes config file content with optional envSpecMap validation.
// If envSpecMap is provided, env() calls in the config will be validated against defined vars.
// workDir is the root config directory used for embed() path resolution.
func GetConfigData(cfgPath string, workDir string, envSpecMap ...EnvSpecMap) (_ []byte, gErr error) {
	var configBytes []byte
	defer result.RecoveryErr(&gErr, func(err error) error {
		// Security: mask config content in logs to prevent sensitive data leakage
		log.Err(err).Str("config_path", cfgPath).Msg("failed to process config data")
		return err
	})

	// If workDir not specified, use the directory of cfgPath
	if workDir == "" {
		workDir = filepath.Dir(cfgPath)
	}

	cfg := &config{workDir: workDir}
	if len(envSpecMap) > 0 && envSpecMap[0] != nil {
		cfg.envSpecMap = envSpecMap[0]
	}

	configBytes = result.Wrap(os.ReadFile(cfgPath)).Expect("failed to read config data: %s", cfgPath)
	configBytes = evalData(configBytes, cfg)
	return configBytes, nil
}

func LoadEnvMap(cfgPath string) EnvSpecMap { return loadEnvConfigMap(cfgPath) }

func loadEnvConfigMap(cfgPath string) EnvSpecMap {
	defer recovery.Exit(func(err error) error {
		log.Err(err).Str("path", cfgPath).Msg("load env config map error")
		return err
	})

	var res Resources
	configBytes := result.Wrap(os.ReadFile(cfgPath)).Expect("failed to read config data: %s", cfgPath)
	assert.Must(yaml.Unmarshal(configBytes, &res), "failed to unmarshal resource config")

	parentDir := filepath.Dir(cfgPath)
	envSpecMap := make(EnvSpecMap)
	// Track which file defined each env var for conflict detection
	envSourceMap := make(map[string]string)

	for _, envPath := range res.PatchEnvs {
		envPath = filepath.Join(parentDir, envPath)
		if pathutil.IsNotExist(envPath) {
			log.Warn().Str("env_path", envPath).Msg("env config path not found")
			continue
		}

		pathList := listAllPath(envPath).Expect("failed to list env config path: %s", envPath)
		for _, p := range pathList {
			if !strings.HasSuffix(p, "."+defaultConfigType) {
				continue
			}

			envConfigBytes := result.Wrap(os.ReadFile(p)).
				Map(bytes.TrimSpace).
				UnwrapOrLog(func(e result.Event) {
					e.Str("env_path", p)
					e.Str(logfields.Msg, "failed to handler env config data")
				})
			if len(envConfigBytes) == 0 {
				continue
			}

			envConfigBytes = evalData(envConfigBytes, &config{workDir: filepath.Dir(cfgPath)})

			// Parse into temporary map to check for conflicts
			var tempEnvMap EnvSpecMap
			result.ErrOf(yaml.Unmarshal(envConfigBytes, &tempEnvMap)).
				MustWithLog(func(e result.Event) {
					// Security: don't log raw env data which may contain secrets
					e.Str("env_path", p)
					e.Str(logfields.Msg, "failed to unmarshal env config")
				})

			// Check for duplicate definitions and merge
			for name, spec := range tempEnvMap {
				name = strings.ToUpper(name)
				if existingSource, exists := envSourceMap[name]; exists {
					log.Panic().
						Str("env_name", name).
						Str("first_defined_in", existingSource).
						Str("redefined_in", p).
						Msg("duplicate env var definition detected in patch_envs, each env var must be defined only once")
				}
				envSourceMap[name] = p
				envSpecMap[name] = spec
			}
		}
	}
	if err := initEnv(envSpecMap); err != nil {
		log.Panic().Err(err).Msg("failed to initialize env")
	}
	return envSpecMap
}

func LoadFromPath[T any](cfgPath string) (*Cfg[T], error) {
	defer recovery.Exit(func(err error) error {
		log.Err(err).Str("config_path", cfgPath).Msg("failed to load config")
		return err
	})

	var val T
	valType := reflect.TypeOf(val)
	for valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	if valType.Kind() != reflect.Struct {
		log.Error().
			Str("config_path", cfgPath).
			Str("type", fmt.Sprintf("%#v", val)).
			Msg("config type not correct")
		return nil, fmt.Errorf("config type not correct")
	}

	envCfgMap := loadEnvConfigMap(cfgPath)
	parentDir := filepath.Dir(cfgPath)

	// Pass envSpecMap to GetConfigData to validate env() calls against defined vars
	// workDir is the root config directory for embed() path resolution
	configBytes := result.Wrap(GetConfigData(cfgPath, parentDir, envCfgMap)).Expect("failed to handler config data")
	defer recovery.Exit(func(err error) error {
		// Security: don't log raw config data which may contain secrets
		log.Err(err).
			Str("config_path", cfgPath).
			Msg("failed to load config")
		return err
	})

	if err := yaml.Unmarshal(configBytes, val); err != nil {
		log.Err(err).
			Str("config_path", cfgPath).
			Msg("failed to unmarshal config")
		return nil, err
	}

	getRealPath := func(pp []string) []string {
		pp = lo.Map(pp, func(item string, index int) string { return filepath.Join(parentDir, item) })

		var resPaths []string
		for _, resPath := range pp {
			pathList := listAllPath(resPath).Expect("failed to list cfgPath: %s", resPath)
			resPaths = append(resPaths, pathList...)
		}

		// skip .*.yaml and cfg.other
		cfgFilter := func(item string, index int) bool {
			return strings.HasSuffix(item, "."+defaultConfigType) && !strings.HasPrefix(item, ".")
		}
		resPaths = lo.Filter(resPaths, cfgFilter)
		return lo.Uniq(resPaths)
	}
	getCfg := func(resPath string) T {
		// Pass envCfgMap to validate env() calls against defined vars
		// Use parentDir as workDir so embed() paths are relative to root config dir
		resBytes := result.Wrap(GetConfigData(resPath, parentDir, envCfgMap)).Expect("failed to handler config data")

		var cfg1 T
		result.ErrOf(yaml.Unmarshal(resBytes, &cfg1)).MustWithLog(func(e result.Event) {
			fmt.Println("res_path", resPath)
			fmt.Println("config_data", string(resBytes))
			assert.Exit(os.WriteFile(resPath+".err.yml", resBytes, 0o666))

			e.Msg("failed to unmarshal config")
		})

		return cfg1
	}

	var res Resources
	assert.Must(yaml.Unmarshal(configBytes, &res), "failed to unmarshal resource config")

	var cfgList []T
	cfgList = append(cfgList, typex.DoBlock1(func() []T {
		resPathList := getRealPath(res.Resources)
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
		patchResPathList := getRealPath(res.PatchResources)
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

	err := Merge(&val, cfgList...)
	if err != nil {
		for _, cfg := range cfgList {
			_, _ = pretty.Simple().Println(cfg)
		}
		log.Fatal().Err(err).Msg("failed to merge config")
	}
	return &Cfg[T]{T: val, P: &val, EnvCfg: lo.ToPtr(envCfgMap)}, nil
}

type Cfg[T any] struct {
	T      T
	P      *T
	EnvCfg *EnvSpecMap
}

// TryLoad attempts to load configuration and returns error instead of panicking.
// This is the recommended way to load config for better error handling.
func TryLoad[T any]() (*Cfg[T], error) {
	var cfgPath = globalManager.GetPath()
	var cfgDir string
	if cfgPath != "" {
		cfgDir = filepath.Dir(cfgPath)
	} else {
		var err error
		cfgPath, cfgDir, err = findConfigPath(defaultConfigName, defaultConfigType)
		if err != nil {
			return nil, err
		}
	}

	globalManager.SetPath(cfgPath)
	globalManager.SetDir(cfgDir)

	return LoadFromPath[T](cfgPath)
}

// Load loads configuration (panics on error).
// For better error handling, use TryLoad instead.
func Load[T any]() Cfg[T] {
	cfg, err := TryLoad[T]()
	assert.Must(err, "failed to load config")
	return lo.FromPtr(cfg)
}
