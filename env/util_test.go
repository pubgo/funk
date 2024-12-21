package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	k, ok := Normalize("aA-bS3_AK/c.d")
	assert.True(t, ok)
	assert.Equal(t, k, "A_A_B_S3_AK_C_D")
}
