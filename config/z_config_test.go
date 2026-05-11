package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/env"
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

func callValidateEnvReferences(cfgPath string) (panicVal any, err error) {
	defer func() {
		panicVal = recover()
	}()
	err = ValidateEnvReferences(cfgPath)
	return
}

func withIsolatedManager(t *testing.T, fn func()) {
	t.Helper()
	oldMgr := globalManager
	globalManager = NewManager()
	t.Cleanup(func() { globalManager = oldMgr })
	fn()
}

func readYAMLAsMap(t *testing.T, p string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(p)
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, yaml.Unmarshal(b, &out))
	return out
}

func resetGoldenRelatedEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"APP_NAME", "MODE", "TIMEOUT"} {
		require.NoError(t, os.Unsetenv(k))
	}
	env.Reload()
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

	panicVal, panicErr := callValidateEnvReferences(cfgPath)
	if panicVal != nil {
		assert.Contains(t, fmt.Sprint(panicVal), "not defined in envs")
		return
	}
	assert.Error(t, panicErr)
	assert.Contains(t, panicErr.Error(), "not defined in envs")
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

func TestComplexFixture_LoadMergedConfigData(t *testing.T) {
	out, err := LoadMergedConfigData("./configs/complex/config.yaml")
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
	assert.Equal(t, "from-patch", app["name"])
	assert.Equal(t, "svc-funk", app["computed_name"])
	assert.Equal(t, "production", app["mode_from_main"])
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("token-123\n")), app["secret_from_embed"])

	database, ok := data["database"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "db.local", database["host"])
	assert.Equal(t, int(6432), database["port"])

	feature, ok := data["feature"].(map[string]any)
	require.True(t, ok)
	tags, ok := feature["tags"].([]any)
	require.True(t, ok)
	assert.Equal(t, []any{"main", "r1", "p1"}, tags)
}

func TestComplexFixture_ValidateEnvReferences_Fails(t *testing.T) {
	panicVal, err := callValidateEnvReferences("./configs/complex_invalid_env/config.yaml")
	if panicVal != nil {
		assert.Contains(t, fmt.Sprint(panicVal), "not defined in envs")
		return
	}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not defined in envs")
}

type goldenCfg struct {
	App struct {
		Name   string   `yaml:"name"`
		Mode   string   `yaml:"mode"`
		Custom string   `yaml:"custom"`
		Labels []string `yaml:"labels"`
	} `yaml:"app"`
	Service struct {
		Retries int    `yaml:"retries"`
		Timeout string `yaml:"timeout"`
	} `yaml:"service"`
}

func TestGolden_LoadMergedConfigData(t *testing.T) {
	withIsolatedManager(t, func() {
		resetGoldenRelatedEnv(t)
		require.NoError(t, RegisterExpr("golden_upper", func(s string) string { return strings.ToUpper(s) }))

		out, err := LoadMergedConfigData("./configs/golden_input_case/config.yaml")
		require.NoError(t, err)

		var got map[string]any
		require.NoError(t, yaml.Unmarshal(out, &got))

		expect := readYAMLAsMap(t, "./configs/golden/merged.golden.yaml")
		assert.Equal(t, expect, got)
	})
}

func TestGolden_LoadFromPath(t *testing.T) {
	withIsolatedManager(t, func() {
		resetGoldenRelatedEnv(t)
		require.NoError(t, RegisterExpr("golden_upper", func(s string) string { return strings.ToUpper(s) }))

		cfg, err := LoadFromPath[goldenCfg]("./configs/golden_input_case/config.yaml")
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "from-patch", cfg.T.App.Name)
		assert.Equal(t, "production", cfg.T.App.Mode)
		assert.Equal(t, "FUNK", cfg.T.App.Custom)
		assert.Equal(t, []string{"base", "patch"}, cfg.T.App.Labels)
		assert.Equal(t, 3, cfg.T.Service.Retries)
		assert.Equal(t, "5s", cfg.T.Service.Timeout)
	})
}

func TestGolden_TryLoadMergedDataAndLoadMergedData(t *testing.T) {
	withIsolatedManager(t, func() {
		resetGoldenRelatedEnv(t)
		require.NoError(t, RegisterExpr("golden_upper", func(s string) string { return strings.ToUpper(s) }))

		oldPath := GetConfigPath()
		defer globalManager.SetPath(oldPath)
		SetConfigPath("./configs/golden_input_case/config.yaml")

		out1, err := TryLoadMergedData()
		require.NoError(t, err)
		out2 := LoadMergedData()

		var got1 map[string]any
		var got2 map[string]any
		require.NoError(t, yaml.Unmarshal(out1, &got1))
		require.NoError(t, yaml.Unmarshal(out2, &got2))

		expect := readYAMLAsMap(t, "./configs/golden/merged.golden.yaml")
		assert.Equal(t, expect, got1)
		assert.Equal(t, expect, got2)
	})
}

func TestGolden_ValidateEnvReferences(t *testing.T) {
	withIsolatedManager(t, func() {
		resetGoldenRelatedEnv(t)
		require.NoError(t, RegisterExpr("golden_upper", func(s string) string { return strings.ToUpper(s) }))

		panicVal, err := callValidateEnvReferences("./configs/golden_input_case/config.yaml")
		if panicVal != nil {
			t.Fatalf("unexpected panic on valid golden config: %v", panicVal)
		}
		assert.NoError(t, err)
	})
}

func TestGolden_ValidateEnvReferences_Fails(t *testing.T) {
	panicVal, err := callValidateEnvReferences("./configs/golden_invalid_input_case/config.yaml")
	if panicVal != nil {
		assert.Contains(t, fmt.Sprint(panicVal), "not defined in envs")
		return
	}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not defined in envs")
}
