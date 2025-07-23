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

func TestEnvPrefix(t *testing.T) {
	Set(PrefixKey, "test")
	Set("test_hello", "world")
	loadEnv()

	envMap := Map()
	assert.Equal(t, envMap["TEST_HELLO"], "world")
}
