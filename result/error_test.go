package result

import (
	"testing"

	"github.com/pubgo/funk/v2/errors"
)

func TestErrorLog(t *testing.T) {
	ErrOf(errors.New("test")).
		Log(func(e Event) {
			e.Msg("ok")
		})
}
