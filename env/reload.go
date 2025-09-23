package env

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pubgo/funk/log/logfields"
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

	logger := slog.With(
		slog.String(logfields.Module, "env"),
		slog.String(logfields.Operation, "reload_env"),
	)
	logger.Info("reload env", slog.Any("env_prefix", map[string]any{"key": PrefixKey, "value": envPrefix}))

	for _, env := range os.Environ() {
		kvs := strings.SplitN(env, "=", 2)
		if len(kvs) != 2 {
			logger.Error("split env error")
			continue
		}

		rawEnvKey := trim(kvs[0])
		if rawEnvKey == "" ||
			strings.HasPrefix(rawEnvKey, "_") ||
			strings.HasPrefix(rawEnvKey, "=") ||
			!hasEnvPrefix(rawEnvKey, envPrefix) {
			logger.Warn(fmt.Sprintf("unset env, key=%s", rawEnvKey))
			_ = os.Unsetenv(rawEnvKey)
			continue
		}

		key, ok := Normalize(rawEnvKey)
		if ok {
			if key == rawEnvKey {
				continue
			}

			setOk := os.Setenv(key, kvs[1]) == nil
			logger.Info(fmt.Sprintf("reset env, old_key=%s new_key=%s set_ok=%v", rawEnvKey, key, setOk))
		} else {
			_ = os.Unsetenv(rawEnvKey)
		}
	}
}
