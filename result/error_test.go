package result

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log/logfields"
)

func TestErrorLog(t *testing.T) {
	ErrOf(errors.New("test")).
		Log(func(e *zerolog.Event) {
			e.Str(logfields.Msg, "ok")
		})
}
