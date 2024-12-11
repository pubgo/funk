package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	k, v, ok := Normalize("aA-b/c.d=1.ok/234-cc_qq")
	assert.True(t, ok)
	assert.Equal(t, k, "A_A_B_C_D")
	assert.Equal(t, v, "1.ok/234-cc_qq")
}
