package env

import (
	"log/slog"
	"strings"

	"github.com/pubgo/funk/v2/internal/strcase"
)

var (
	trim     = strings.TrimSpace
	replacer = strcase.NewCaser(
		true,
		map[string]bool{"SSL": true, "HTML": false},
		strcase.NewSplitFn(
			[]rune{'*', '.', ',', '-', '/'},
			strcase.SplitCase,
			strcase.SplitAcronym,
			strcase.PreserveNumberFormatting,
		))
)

func keyHandler(key string) string {
	key = strings.ReplaceAll(replacer.ToSNAKE(key), "__", "_")
	return strings.ToUpper(trim(key))
}

// getLog resolves slog.Default on every call (not once) so env logs follow
// the application's current default logger, even when slog.SetDefault happens
// after the first env log. WithGroup allocates, but env logs rarely.
func getLog() *slog.Logger {
	return slog.Default().WithGroup(Name)
}
