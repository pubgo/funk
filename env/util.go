package env

import (
	"strings"
)

var trim = strings.TrimSpace
var replacer = strings.NewReplacer("-", "_", ".", "_", "/", "_")

func KeyHandler(key string) string {
	return strings.ToUpper(trim(strings.ReplaceAll(replacer.Replace(key), "__", "_")))
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return KeyHandler(key), true
}
