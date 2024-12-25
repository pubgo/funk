package env

import (
	"strings"

	strcase "github.com/ettle/strcase"
)

var replacer = strcase.NewCaser(
	true,
	map[string]bool{"SSL": true, "HTML": false},
	strcase.NewSplitFn(
		[]rune{'*', '.', ',', '-', '/'},
		strcase.SplitCase,
		strcase.SplitAcronym,
		strcase.PreserveNumberFormatting,
	))
var trim = strings.TrimSpace

func KeyHandler(key string) string {
	return strings.ToUpper(trim(strings.ReplaceAll(replacer.ToSNAKE(key), "__", "_")))
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return KeyHandler(key), true
}
