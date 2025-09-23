package env

import (
	"os"
	"strings"

	"github.com/ettle/strcase"
)

const PrefixKey = "ENV_PREFIX"

func hasEnvPrefix(key string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(prefix))
}

func getEnvPrefix() string {
	prefix := strings.TrimSpace(os.Getenv(PrefixKey))
	if prefix != "" {
		prefix = strings.ReplaceAll(prefix+"_", "__", "_")
	}
	return strings.ToUpper(prefix)
}

var trim = strings.TrimSpace
var replacer = strcase.NewCaser(
	true,
	map[string]bool{"SSL": true, "HTML": false},
	strcase.NewSplitFn(
		[]rune{'*', '.', ',', '-', '/'},
		strcase.SplitCase,
		strcase.SplitAcronym,
		strcase.PreserveNumberFormatting,
	))

func KeyHandler(key string) string {
	key = strings.ToUpper(replacer.ToSNAKE(key))
	envPrefix := getEnvPrefix()
	if envPrefix != "" {
		key = envPrefix + "_" + strings.TrimPrefix(key, envPrefix)
	}
	return strings.ToUpper(trim(strings.ReplaceAll(key, "__", "_")))
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return KeyHandler(key), true
}
