package assert_test

import (
	"testing"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
)

func TestCheckNil(t *testing.T) {
	var a *int

	defer recovery.DebugPrint()

	assert.Assert(a == nil, "ok")
}
