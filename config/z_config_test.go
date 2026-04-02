package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEnvMap(t *testing.T) {
	envs := LoadEnvMap("./configs/config.yaml")
	assert.NotNil(t, envs["TEST1"])
	assert.NotNil(t, envs["TEST2"])
}

type envRefTestCfg struct {
	App struct {
		Declared   string `yaml:"declared"`
		Undeclared string `yaml:"undeclared"`
	} `yaml:"app"`
}

func TestLoadFromPath_DoesNotForceEnvReferenceValidation(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
DECLARED_VAR:
  default: declared
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "app.yaml"), []byte(`
app:
  declared: ${DECLARED_VAR:default_declared}
  undeclared: ${UNDECLARED_VAR:default_undeclared}
`), 0o644))

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - app.yaml
patch_envs:
  - envs
`), 0o644))

	_, err := LoadFromPath[envRefTestCfg](cfgPath)
	assert.NoError(t, err)
}

func TestValidateEnvReferences_FailsOnUndeclaredEnv(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
DECLARED_VAR:
  default: declared
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "app.yaml"), []byte(`
app:
  declared: ${DECLARED_VAR:default_declared}
  undeclared: ${UNDECLARED_VAR:default_undeclared}
`), 0o644))

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - app.yaml
patch_envs:
  - envs
`), 0o644))

	err := ValidateEnvReferences(cfgPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not defined in envs")
}

func TestValidateEnvReferences_Success(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
DECLARED_VAR:
  default: declared
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "app.yaml"), []byte(`
app:
  declared: ${DECLARED_VAR:default_declared}
`), 0o644))

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - app.yaml
patch_envs:
  - envs
`), 0o644))

	err := ValidateEnvReferences(cfgPath)
	assert.NoError(t, err)
}

func TestLoadMergedConfigData(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "envs"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "envs", "base.yaml"), []byte(`
DECLARED:
  default: merged_from_patch_envs
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "base.yaml"), []byte(`
app:
  name: base
  from_base: true
  env_from_base: ${DECLARED:base_default}
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "override.yaml"), []byte(`
app:
  name: override
  from_patch: true
`), 0o644))

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
resources:
  - base.yaml
patch_resources:
  - override.yaml
patch_envs:
  - envs
app:
  from_main: true
  env_from_main: ${DECLARED:main_default}
`), 0o644))

	out, err := LoadMergedConfigData(cfgPath)
	require.NoError(t, err)
	require.NotEmpty(t, out)

	var data map[string]any
	require.NoError(t, yaml.Unmarshal(out, &data))

	_, hasResources := data["resources"]
	_, hasPatchResources := data["patch_resources"]
	_, hasPatchEnvs := data["patch_envs"]
	assert.False(t, hasResources)
	assert.False(t, hasPatchResources)
	assert.False(t, hasPatchEnvs)

	app, ok := data["app"].(map[string]any)
	require.True(t, ok)

	assert.Equal(t, "override", app["name"])
	assert.Equal(t, true, app["from_main"])
	assert.Equal(t, true, app["from_base"])
	assert.Equal(t, true, app["from_patch"])
	assert.Equal(t, "merged_from_patch_envs", app["env_from_main"])
	assert.Equal(t, "merged_from_patch_envs", app["env_from_base"])
}

func TestLoadMergedConfigData_EmptyPath(t *testing.T) {
	_, err := LoadMergedConfigData("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config path is null")
}

func TestTryLoadMergedData(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	oldPath := GetConfigPath()
	defer globalManager.SetPath(oldPath)

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
app:
  name: demo
`), 0o644))

	SetConfigPath(cfgPath)
	data, err := TryLoadMergedData()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var got map[string]any
	require.NoError(t, yaml.Unmarshal(data, &got))
	app, ok := got["app"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "demo", app["name"])
}

func TestLoadMergedData(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	oldPath := GetConfigPath()
	defer globalManager.SetPath(oldPath)

	require.NoError(t, os.WriteFile(cfgPath, []byte(`
app:
  enabled: true
`), 0o644))

	SetConfigPath(cfgPath)
	data := LoadMergedData()
	require.NotEmpty(t, data)

	var got map[string]any
	require.NoError(t, yaml.Unmarshal(data, &got))
	app, ok := got["app"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, app["enabled"])
}
