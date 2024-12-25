package env

import (
	"os"
	"strings"
)

func Init() {
	initEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func initEnv() {
	for _, env := range os.Environ() {
		kvs := strings.SplitN(env, "=", 2)
		if len(kvs) != 2 {
			continue
		}

		var rawKey = kvs[0]
		key, ok := Normalize(rawKey)
		if !ok {
			_ = os.Unsetenv(rawKey)
			continue
		}

		_ = os.Unsetenv(rawKey)
		_ = os.Setenv(key, kvs[1])
	}
}

func init() {
	initEnv()
}
