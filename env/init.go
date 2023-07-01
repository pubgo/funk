package env

import (
	"os"
)

// 环境变量处理, key转大写, 同时把`-./`转换为`_`
// a-b=>a_b, a.b=>a_b, a/b=>a_b
func init() {
	for _, env := range os.Environ() {
		k, v, ok := Normalize(env)
		if !ok {
			if k != "" {
				_ = os.Unsetenv(k)
			}
		} else {
			_ = os.Setenv(k, v)
		}
	}
}
