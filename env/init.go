package env

import (
	"os"
	"strings"
)

var Prefix = ""

func GetPrefix() string {
	return Key(Prefix)
}

func Init() {
	initEnv()
}

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func initEnv() {
	for _, env := range os.Environ() {
		k, v, ok := Normalize(env)
		if k != "" && ok && strings.HasPrefix(k, GetPrefix()) {
			_ = os.Setenv(k, v)
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
