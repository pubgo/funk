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

	// MustDelete propagates the same cleanup.
	env.MustSet("DEL_MUST", "v")
	env.MustDelete("DEL_MUST")
	_, ok = env.Lookup("DEL_MUST")
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
	setExpandEnv(t, "GET_INT_PAD", " 42 ")

	// Get/getters trim surrounding whitespace, matching the v2 behavior.
	assert.Equal(t, "42", env.Get("GET_INT_PAD"))
	assert.Equal(t, 42, env.GetInt("GET_INT_PAD"))

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
	assert.Equal(t, 0.5, env.GetFloatOr("GET_INT_BAD", 0.5))
	assert.Equal(t, 1.5, env.GetFloatOr("GET_FLOAT", 0.5))

	// Wrapper APIs delegate to the covered paths.
	setExpandEnv(t, "GET_MUST", "v")
	assert.Equal(t, "v", env.MustGet("GET_MUST"))
	assert.Panics(t, func() { env.MustGet("GET_UNSET") })

	var got string
	env.GetWith(&got, "GET_MUST", "GET_UNSET")
	assert.Equal(t, "v", got)
	env.GetWith(&got, "GET_UNSET", "GET_MUST")
	assert.Equal(t, "v", got, "first non-empty candidate wins")
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

func TestMustLoadFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "must.env")
	require.NoError(t, os.WriteFile(path, []byte("DOTENV_MUST=ok\n"), 0o644))
	t.Cleanup(func() { require.NoError(t, os.Unsetenv("DOTENV_MUST")) })

	assert.NotPanics(t, func() { env.MustLoadFiles(path) })
	assert.Equal(t, "ok", env.Get("DOTENV_MUST"))

	// A directory passes the existence check but fails ReadFile -> panic.
	assert.Panics(t, func() { env.MustLoadFiles(t.TempDir()) })
}

func TestDeleteRemovesAllSpellings(t *testing.T) {
	// env.MustSet stores under the normalized key; a raw export may also
	// exist. Delete must unset every candidate spelling.
	env.MustSet("T9KEY", "via-set")
	require.NoError(t, os.Setenv("T9KEY", "raw-value"))
	t.Cleanup(func() {
		require.NoError(t, os.Unsetenv("T9KEY"))
		require.NoError(t, os.Unsetenv("T9_KEY"))
	})

	_, ok := env.Lookup("T9_KEY")
	require.True(t, ok)

	assert.NoError(t, env.Delete("T9KEY").Err())
	_, ok = env.Lookup("T9KEY")
	assert.False(t, ok, "raw spelling should be unset")
	_, ok = env.Lookup("T9_KEY")
	assert.False(t, ok, "normalized spelling should be unset")
}

func TestLookupEmptyValue(t *testing.T) {
	// Lookup reports existence even for empty values; Get skips them.
	setExpandEnv(t, "LOOKUP_EMPTY", "")

	v, ok := env.Lookup("LOOKUP_EMPTY")
	assert.True(t, ok, "set-but-empty should be reported as present")
	assert.Equal(t, "", v)
	assert.Equal(t, "", env.Get("LOOKUP_EMPTY"))
	assert.Equal(t, "def", env.GetOr("LOOKUP_EMPTY", "def"))
}

func TestGetOrTrimsValue(t *testing.T) {
	setExpandEnv(t, "GET_OR_PAD", " padded ")

	assert.Equal(t, "padded", env.Get("GET_OR_PAD"))
	assert.Equal(t, "padded", env.GetOr("GET_OR_PAD", "default"))
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
