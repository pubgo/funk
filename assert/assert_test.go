package assert

import (
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestCheckNil(t *testing.T) {
	var a *int

	defer func() {
		errors.Debug(errors.Parse(recover()))
	}()

	Assert(a == nil, "ok")
}
