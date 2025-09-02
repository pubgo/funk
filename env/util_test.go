package env_test

import (
	"os"
	"strings"
	"testing"

	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/pretty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	k, ok := env.Normalize("aA-bS3_AK/c.d")
	assert.True(t, ok)
	assert.Equal(t, k, "A_A_B_S3_AK_C_D")
}

func TestEnvPrefix(t *testing.T) {
	log.Logger = log.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		if strings.HasPrefix(message, "unset not match env") {
			e.Discard()
		}
	}))

	env.Set(env.PrefixKey, "test")
	env.Set("test_hello", "world")
	env.Reload()

	envMap := env.Map()
	assert.Equal(t, envMap["TEST_HELLO"], "world")
	pretty.Println(os.Environ())
}
