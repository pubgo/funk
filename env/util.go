package env

import (
	"strings"
)

var replacer = strings.NewReplacer("-", "_", ".", "_", "/", "_")

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(env string) (k, v string, ok bool) {
	if env == "" {
		return "", "", false
	}

	envs := strings.SplitN(env, "=", 2)
	var key = trim(envs[0])
	if len(envs) != 2 || key == "" || strings.HasPrefix(key, "_") {
		return key, "", false
	}

	key = replacer.Replace(key)
	key = strings.ToUpper(strings.ReplaceAll(key, "__", "_"))
	return key, trim(envs[1]), true
}
