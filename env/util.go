package env

import (
	"strings"
)

var replacer = strings.NewReplacer("-", "_", ".", "_", "/", "_")

func KeyHandler(key string) string {
	return strings.ToUpper(trim(strings.ReplaceAll(replacer.Replace(key), "__", "_")))
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(env string) (k, v string, ok bool) {
	if env == "" {
		return "", "", false
	}

	envs := strings.SplitN(env, "=", 2)
	key := trim(envs[0])
	if len(envs) != 2 || key == "" || strings.HasPrefix(key, "_") {
		return key, "", false
	}

	return KeyHandler(key), trim(envs[1]), true
}
