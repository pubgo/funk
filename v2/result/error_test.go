package result

import (
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestErrorLog(t *testing.T) {
	ErrOf(errors.New("test")).LogErr()
}
