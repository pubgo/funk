package env_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/log"
)

func init() {
	slog.SetDefault(slog.New(log.NewSlog(log.GetLogger(""))))
}

func TestResetEnv(t *testing.T) {
	assert.NoError(t, os.Setenv("abc", "1"))
	assert.Equal(t, os.Getenv("abc"), "1")
	assert.NoError(t, os.Setenv("abc", "2"))
	assert.Equal(t, os.Getenv("abc"), "2")
}

func TestNormalize(t *testing.T) {
	k, ok := env.Normalize("aA-bS3_AK/c.d")
	assert.True(t, ok)
	assert.Equal(t, k, "A_A_B_S3_AK_C_D")
}

func TestEnvPrefix(t *testing.T) {
	env.Reload()

	env.MustSet("test_hello", "world")
	env.Reload()

	envMap := env.Map()
	assert.Equal(t, envMap["TEST_HELLO"], "world")
}
