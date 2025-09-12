package assert_test

import (
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
)

func TestCheckNil(t *testing.T) {
	var a *int

	defer recovery.DebugPrint()

	assert.Assert(a == nil, "ok")
}
