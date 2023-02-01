package try

import (
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestTry(t *testing.T) {
	errors.Debug(Try(func() error {
		panic("hello")
	}))
}
