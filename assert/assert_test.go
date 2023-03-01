package assert

import (
	"github.com/pubgo/funk/errors"
	"testing"
)

func TestCheckNil(t *testing.T) {
	var a *int

	defer func() {
		errors.Debug(errors.Parse(recover()))
	}()

	Assert(a == nil, "ok")
}
