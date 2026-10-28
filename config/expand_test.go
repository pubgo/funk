package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type envSubstCfg struct {
	App struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
		Sha  string `yaml:"sha"`
		Mode string `yaml:"mode"`
	} `yaml:"app"`
}

func unsetEnvs(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		require.NoError(t, os.Unsetenv(k))
	}
}

// ${PORT:8080} with PORT unset must resolve to the default, not silently drop it.
func TestEnvSubstNumericDefault(t *testing.T) {
	withIsolatedManager(t, func() {
		unsetEnvs(t, "PORT", "HOST")

		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(`
app:
  port: ${PORT:8080}
  host: ${HOST:localhost}
`), 0o644))

		cfg, err := LoadFromPath[envSubstCfg](cfgPath)
		assert.NoError(t, err)
		assert.Equal(t, 8080, cfg.T.App.Port)
		assert.Equal(t, "localhost", cfg.T.App.Host)
	})
}

// Keys containing digits get renamed by keyHandler on env.Set (T2HOST -> T2_HOST);
// substitution must still resolve them through the original name.
func TestEnvSubstDigitKey(t *testing.T) {
	withIsolatedManager(t, func() {
		unsetEnvs(t, "T2HOST", "T2_HOST")
		t.Cleanup(func() { unsetEnvs(t, "T2_HOST") })

		tmpDir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
T2HOST:
  default: t2.example.com
SHA256SUM:
  default: deadbeef
`), 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "app.yaml"), []byte(`
app:
  host: ${T2HOST:fallback}
  sha: ${SHA256SUM:fallback}
`), 0o644))
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - app.yaml
patch_envs:
  - envs
`), 0o644))

		cfg, err := LoadFromPath[envSubstCfg](cfgPath)
		assert.NoError(t, err)
		assert.Equal(t, "t2.example.com", cfg.T.App.Host)
		assert.Equal(t, "deadbeef", cfg.T.App.Sha)
	})
}

// A raw process env var with a digit key, never declared in patch_envs, must also resolve.
func TestEnvSubstRawDigitKey(t *testing.T) {
	withIsolatedManager(t, func() {
		require.NoError(t, os.Setenv("V2RAY", "v2-value"))
		t.Cleanup(func() { unsetEnvs(t, "V2RAY") })

		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(`
app:
  host: ${V2RAY:fallback}
`), 0o644))

		cfg, err := LoadFromPath[envSubstCfg](cfgPath)
		assert.NoError(t, err)
		assert.Equal(t, "v2-value", cfg.T.App.Host)
	})
}

// Declared patch_envs values keep priority over inline defaults (existing behavior).
func TestEnvSubstDeclaredDefaultWins(t *testing.T) {
	withIsolatedManager(t, func() {
		unsetEnvs(t, "MODE")
		t.Cleanup(func() { unsetEnvs(t, "MODE") })

		tmpDir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
MODE:
  default: production
`), 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "app.yaml"), []byte(`
app:
  mode: ${MODE:dev}
`), 0o644))
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - app.yaml
patch_envs:
  - envs
`), 0o644))

		cfg, err := LoadFromPath[envSubstCfg](cfgPath)
		assert.NoError(t, err)
		assert.Equal(t, "production", cfg.T.App.Mode)
	})
}
