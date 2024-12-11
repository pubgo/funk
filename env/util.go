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
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return KeyHandler(key), true
}
