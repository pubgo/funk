package env

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/ettle/strcase"
)

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

func keyHandler(key string) string {
	key = strings.ReplaceAll(replacer.ToSNAKE(key), "__", "_")
	return strings.ToUpper(trim(key))
}

var getLog = sync.OnceValue(func() *slog.Logger {
	return slog.Default().WithGroup(Name)
})
