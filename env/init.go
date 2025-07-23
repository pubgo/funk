package env

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	loadEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func loadEnv() {
	envPrefix := getEnvPrefix()
	envPrefixEventFn := func(e *zerolog.Event) {
		e.Dict("env_prefix", zerolog.Dict().Str("key", PrefixKey).Str("value", envPrefix))
	}
	logRecord(log.Info(), envPrefixEventFn).Msg("load env")

	for _, env := range os.Environ() {
		rawEnvFn := func(e *zerolog.Event) { e.Str("raw_env", env) }
		kvs := strings.SplitN(env, "=", 2)
		if len(kvs) != 2 {
			logRecord(log.Error(), envPrefixEventFn, rawEnvFn).Msg("split env error")
			continue
		}

		envKey := trim(kvs[0])
		_ = os.Unsetenv(envKey)

		if envKey == "" ||
			strings.HasPrefix(envKey, "_") ||
			strings.HasPrefix(envKey, "=") ||
			!hasEnvPrefix(envKey, envPrefix) {
			logRecord(log.Warn(), envPrefixEventFn, rawEnvFn).Msgf("ignore env, key=%s", envKey)
			continue
		}

		key, ok := Normalize(envKey)
		if ok {
			_ = os.Setenv(key, kvs[1])
		}

		logRecord(log.Info(), envPrefixEventFn, rawEnvFn).Msgf("reset env, old_env_key=%s new_env_key=%s", envKey, key)
	}
}
