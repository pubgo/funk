package env

import (
	"os"
)

func Init() {
	initEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func initEnv() {
	for _, env := range os.Environ() {
		k, v, ok := Normalize(env)
		if k != "" && ok {
			_ = os.Setenv(k, v)
			continue
		}

		if k == "" {
			continue
		}

		_ = os.Unsetenv(k)
	}
}

func init() {
	initEnv()
}
