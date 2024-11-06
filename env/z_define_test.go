package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefine(t *testing.T) {
	var a = Define("env_test", "desc_a")
	var b = Define("env_test", "desc_b")
	assert.NoError(t, Set("env_test", "a"))
	assert.Equal(t, a.Get(), "a")
	assert.NoError(t, Set("env_test", "b"))
	assert.Equal(t, b.Get(), "b")
	for _, ee := range FindAllDefinedEnvs() {
		assert.Equal(t, ee.Key, "ENV_TEST")
	}
}
