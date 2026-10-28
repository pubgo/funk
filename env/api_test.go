package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgo/funk/v2/env"
)

func TestGetLookupDeleteRawFallback(t *testing.T) {
	// Raw export, never passed through Set/Reload.
	require.NoError(t, os.Setenv("V2RAY", "raw-value"))
	t.Cleanup(func() { require.NoError(t, os.Unsetenv("V2RAY")) })

	assert.Equal(t, "raw-value", env.Get("V2RAY"))

	v, ok := env.Lookup("V2RAY")
	assert.True(t, ok)
	assert.Equal(t, "raw-value", v)

	assert.NoError(t, env.Delete("V2RAY").Err())
	_, ok = env.Lookup("V2RAY")
	assert.False(t, ok)
}

func TestGetRenamedKey(t *testing.T) {
	// Set stores under the normalized key (T2HOST -> T2_HOST); all
	// spellings must resolve, matching env.Expand behavior.
	env.MustSet("T2HOST", "t2-value")
	t.Cleanup(func() { require.NoError(t, os.Unsetenv("T2_HOST")) })

	assert.Equal(t, "t2-value", env.Get("T2HOST"))
	assert.Equal(t, "t2-value", env.Get("t2host"))
	assert.Equal(t, "t2-value", env.Get("T2_HOST"))

	v, ok := env.Lookup("t2host")
	assert.True(t, ok)
	assert.Equal(t, "t2-value", v)
}

func TestTypedGetters(t *testing.T) {
	setExpandEnv(t, "GET_BOOL", "true")
	setExpandEnv(t, "GET_INT", "42")
	setExpandEnv(t, "GET_FLOAT", "1.5")
	setExpandEnv(t, "GET_INT_BAD", "abc")

	assert.True(t, env.GetBool("GET_BOOL"))
	assert.Equal(t, 42, env.GetInt("GET_INT"))
	assert.Equal(t, 1.5, env.GetFloat("GET_FLOAT"))

	// Unset -> zero sentinels.
	assert.False(t, env.GetBool("GET_UNSET"))
	assert.Equal(t, -1, env.GetInt("GET_UNSET"))
	assert.Equal(t, -1.0, env.GetFloat("GET_UNSET"))

	// Invalid value -> zero sentinels.
	assert.False(t, env.GetBool("GET_INT_BAD"))
	assert.Equal(t, -1, env.GetInt("GET_INT_BAD"))

	// Or variants fall back on unset and invalid values.
	assert.True(t, env.GetBoolOr("GET_UNSET", true))
	assert.False(t, env.GetBoolOr("GET_INT_BAD", false))
	assert.Equal(t, 7, env.GetIntOr("GET_UNSET", 7))
	assert.Equal(t, 7, env.GetIntOr("GET_INT_BAD", 7))
	assert.Equal(t, 42, env.GetIntOr("GET_INT", 7))
	assert.Equal(t, 0.5, env.GetFloatOr("GET_UNSET", 0.5))
}

func TestLoadFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.env")
	require.NoError(t, os.WriteFile(path, []byte(`
DOTENV_A=hello
DOTENV_B=42
# comment and empty values are skipped
DOTENV_EMPTY=
`), 0o644))

	assert.NoError(t, env.LoadFiles(path).Err())
	t.Cleanup(func() {
		require.NoError(t, os.Unsetenv("DOTENV_A"))
		require.NoError(t, os.Unsetenv("DOTENV_B"))
	})

	assert.Equal(t, "hello", env.Get("DOTENV_A"))
	assert.Equal(t, 42, env.GetInt("DOTENV_B"))
	assert.Equal(t, "", env.Get("DOTENV_EMPTY"))

	// Missing files are ignored, no error.
	assert.NoError(t, env.LoadFiles(filepath.Join(t.TempDir(), "missing.env")).Err())
}

func TestExpandAfterReload(t *testing.T) {
	// A raw export renamed by Reload must still resolve via the
	// normalized candidate.
	require.NoError(t, os.Setenv("V2RAY", "reloaded"))
	env.Reload()
	t.Cleanup(func() { require.NoError(t, os.Unsetenv("V2_RAY")) })

	assert.Equal(t, "reloaded", env.Expand("${V2RAY}").UnwrapOrEmpty())
	assert.Equal(t, "reloaded", env.Get("V2RAY"))
}
