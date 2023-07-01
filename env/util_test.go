package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	k, v, ok := Normalize("a-b/c.d=1")
	assert.True(t, ok)
	assert.Equal(t, k, "A_B_C_D")
	assert.Equal(t, v, "1")
}
