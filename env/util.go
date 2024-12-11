package env

import (
	"strings"

	"github.com/iancoleman/strcase"
)

var trim = strings.TrimSpace
var replacer = strings.NewReplacer("-", "_", ".", "_", "/", "_")

func KeyHandler(key string) string {
	key = strcase.ToScreamingSnake(key)
	return strings.ToUpper(trim(strings.ReplaceAll(replacer.Replace(key), "__", "_")))
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(env string) (k, v string, ok bool) {
	if env == "" {
		return "", "", false
	}

	kvs := strings.SplitN(env, "=", 2)
	key := trim(kvs[0])
	if len(kvs) != 2 || key == "" || strings.HasPrefix(key, "_") {
		return key, "", false
	}

	return KeyHandler(key), trim(kvs[1]), true
}
