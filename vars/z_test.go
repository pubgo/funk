package vars_test

import (
	"expvar"
	"testing"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/vars"
)

func TestAny(t *testing.T) {
	defer recovery.Testing(t)
	name := "test-bool-var"
	bb := vars.Bool(name)
	assert.If(bb.Load() != false, "not match")
	assert.If(vars.Bool(name).String() != "false", "not match")
	vars.Bool(name).Store(true)
	assert.MustEqual(expvar.Get(name).String(), "true")
}
