package env

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Reload() {
	loadEnv()
}

// Init reload env
// Deprecated: use Reload instead.
func Init() {
	loadEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func loadEnv() {
	envPrefix := getEnvPrefix()
	logger := log.With().Str("operation", "reload_env").Logger()
	envPrefixEventFn := func(e *zerolog.Event) {
		e.Dict("env_prefix", zerolog.Dict().Str("key", PrefixKey).Str("value", envPrefix))
	}
	logRecord(logger.Info(), envPrefixEventFn).Msg("reload env")

	for _, env := range os.Environ() {
		kvs := strings.SplitN(env, "=", 2)
		if len(kvs) != 2 {
			logRecord(logger.Error()).Msg("split env error")
			continue
		}

		rawEnvKey := trim(kvs[0])
		if rawEnvKey == "" ||
			strings.HasPrefix(rawEnvKey, "_") ||
			strings.HasPrefix(rawEnvKey, "=") ||
			!hasEnvPrefix(rawEnvKey, envPrefix) {
			logRecord(logger.Warn()).Msgf("unset env, key=%s", rawEnvKey)
			_ = os.Unsetenv(rawEnvKey)
			continue
		}

		key, ok := Normalize(rawEnvKey)
		if ok {
			if key == rawEnvKey {
				continue
			}

			setOk := os.Setenv(key, kvs[1]) == nil
			logRecord(logger.Info()).Msgf("reset env, old_key=%s new_key=%s set_ok=%v", rawEnvKey, key, setOk)
		} else {
			_ = os.Unsetenv(rawEnvKey)
		}
	}
}
