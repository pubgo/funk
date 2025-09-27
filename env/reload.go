package env

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/log/logfields"
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
	logger := slog.With(slog.String(logfields.Module, "env"), slog.String(logfields.Operation, "reload_env"))
	logger.Info("reload env")

	for _, env := range os.Environ() {
		kvs := strings.SplitN(env, "=", 2)
		if len(kvs) != 2 {
			logger.Error("split env error", slog.String("env", env))
			continue
		}

		oldKey := trim(kvs[0])
		newKey, ok := Normalize(oldKey)
		if ok {
			if newKey == oldKey {
				continue
			}

			assert.Exit(os.Unsetenv(oldKey))
			assert.Exit(os.Setenv(newKey, kvs[1]))
			logger.Info(fmt.Sprintf("reset env, old_key=%s new_key=%s", oldKey, newKey))
		} else {
			logger.Warn(fmt.Sprintf("unset env, key=%s", oldKey))
			assert.Exit(os.Unsetenv(oldKey))
		}
	}
}
