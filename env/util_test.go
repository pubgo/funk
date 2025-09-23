package env_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/pretty"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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
	slog.SetDefault(slog.New(log.NewSlog(log.GetLogger(""))))

	env.Reload()
	pretty.Println("env_keys", lo.Keys(env.Map()))

	env.Set(env.PrefixKey, "test").Must()
	env.Set("test_hello", "world").Must()
	env.Reload()

	envMap := env.Map()
	assert.Equal(t, envMap["TEST_HELLO"], "world")
	pretty.Println(os.Environ())
}
