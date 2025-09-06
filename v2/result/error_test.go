package result

import (
	"github.com/rs/zerolog"
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestErrorLog(t *testing.T) {
	ErrOf(errors.New("test")).
		Log().
		LogCtx(nil, func(e *zerolog.Event) {
			e.Str("abc", "test")
		})
}
