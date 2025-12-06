package try

import (
	"testing"

	"github.com/pubgo/funk/v2/errors"
)

func TestTry(t *testing.T) {
	errors.DebugPrint(Try(func() error {
		panic("hello")
	}))
}
