package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/env"
)

func setExpandEnv(t *testing.T, key, val string) {
	t.Helper()
	assert.NoError(t, os.Setenv(key, val))
	t.Cleanup(func() { assert.NoError(t, os.Unsetenv(key)) })
}

func TestExpand(t *testing.T) {
	// Keys with a digit->letter boundary get renamed by keyHandler on Set.
	assert.Equal(t, "T2_HOST", env.Key("T2HOST"))
	assert.Equal(t, "SHA256_SUM", env.Key("SHA256SUM"))

	env.MustSet("T2HOST", "t2.example.com")
	t.Cleanup(func() { assert.NoError(t, os.Unsetenv("T2_HOST")) })

	env.MustSet("SHA256SUM", "abc123")
	t.Cleanup(func() { assert.NoError(t, os.Unsetenv("SHA256_SUM")) })

	setExpandEnv(t, "EXPAND_RAW", "raw_value")
	setExpandEnv(t, "EXPAND_EMPTY", "")
	setExpandEnv(t, "EXPAND_URL", "http://x:1")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no placeholder", "plain text", "plain text"},
		{"raw key", "${EXPAND_RAW}", "raw_value"},
		{"bare form", "$EXPAND_RAW", "raw_value"},
		{"renamed key via braced form", "${T2HOST}", "t2.example.com"},
		{"renamed key via bare form", "$T2HOST", "t2.example.com"},
		{"normalized key direct", "${T2_HOST}", "t2.example.com"},
		{"sha key", "${SHA256SUM}", "abc123"},
		{"sha key lowercase reference", "${sha256sum}", "abc123"},
		{"sha key mixed-case reference", "${Sha256Sum}", "abc123"},
		{"unset empty", "${EXPAND_UNSET}", ""},
		{"colon default", "${EXPAND_UNSET:8080}", "8080"},
		{"dash default", "${EXPAND_UNSET:-8080}", "8080"},
		{"quoted default", `${EXPAND_UNSET:"localhost"}`, "localhost"},
		{"single-quoted default", "${EXPAND_UNSET:'localhost'}", "localhost"},
		{"set beats default", "${EXPAND_RAW:8080}", "raw_value"},
		{"empty value falls back to default", "${EXPAND_EMPTY:fallback}", "fallback"},
		{"default keeps colons", "${EXPAND_UNSET:http://localhost:8080}", "http://localhost:8080"},
		{"url value", "${EXPAND_URL:http://fallback:2}", "http://x:1"},
		{"multiple placeholders", "port=${EXPAND_UNSET:8080} host=${EXPAND_RAW}", "port=8080 host=raw_value"},
		{"default after first colon", "${EXPAND_UNSET:0:3}", "0:3"},
		{"digit-leading bare name is literal", "cost $5 and $1x", "cost $5 and $1x"},
		{"dollar escape", "a$$b", "a$b"},
		{"lone dollar", "$", "$"},
		{"unterminated braces literal", "${EXPAND_UNSET", "${EXPAND_UNSET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, env.Expand(tt.input).UnwrapOrEmpty())
		})
	}
}
