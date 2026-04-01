package env

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pubgo/funk/v2/assert"
)

func Reload() {
	loadEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func loadEnv() {
	logger := getLog().With(slog.String("logger", Name))
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
			logger.Debug(fmt.Sprintf("reset env, old_key=%s new_key=%s", oldKey, newKey))
		} else {
			logger.Debug(fmt.Sprintf("unset env, key=%s", oldKey))
			assert.Exit(os.Unsetenv(oldKey))
		}
	}
}
