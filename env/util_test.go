package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	k, ok := Normalize("aA-bS3_AK/c.d")
	assert.True(t, ok)
	assert.Equal(t, k, "AA_BS3_AK_C_D")
}
