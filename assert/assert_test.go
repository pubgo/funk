package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckNil(t *testing.T) {
	var is = assert.New(t)
	var a *int

	is.Panics(func() {
		Assert(a == nil, "ok")
	})
}
