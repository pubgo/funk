package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/env"
)

func TestCelEngine_BasicExpressions(t *testing.T) {
	cfg := &config{workDir: t.TempDir()}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		expr     string
		expected any
	}{
		{
			name:     "string literal",
			expr:     `"hello"`,
			expected: "hello",
		},
		{
			name:     "integer literal",
			expr:     `42`,
			expected: int64(42),
		},
		{
			name:     "boolean literal",
			expr:     `true`,
			expected: true,
		},
		{
			name:     "string concatenation",
			expr:     `"hello" + " " + "world"`,
			expected: "hello world",
		},
		{
			name:     "arithmetic",
			expr:     `1 + 2 * 3`,
			expected: int64(7),
		},
		{
			name:     "config_dir function",
			expr:     `config_dir()`,
			expected: cfg.workDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Eval(tt.expr)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCelEngine_EnvAccess(t *testing.T) {
	assert.NoError(t, os.Setenv("CEL_TEST_VAR", "test_value"))
	env.Reload()
	defer func() {
		assert.NoError(t, os.Unsetenv("CEL_TEST_VAR"))
	}()

	// Test env function with envSpecMap - variable must be defined
	cfg := &config{
		workDir: t.TempDir(),
		envSpecMap: EnvSpecMap{
			"CEL_TEST_VAR": &EnvSpec{Name: "CEL_TEST_VAR"},
		},
	}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	// Test env function - should work when var is defined in envSpecMap
	result, err := engine.Eval(`env("CEL_TEST_VAR")`)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", result)

	// Test env with undefined var - should error
	_, err = engine.Eval(`env("NON_EXISTENT_VAR")`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not defined in patch_envs")
}

func TestCelEngine_EnvsFunction(t *testing.T) {
	assert.NoError(t, os.Setenv("ENV_A", "value_a"))
	assert.NoError(t, os.Setenv("ENV_B", "value_b"))
	env.Reload()
	defer func() {
		assert.NoError(t, os.Unsetenv("ENV_A"))
		assert.NoError(t, os.Unsetenv("ENV_B"))
	}()

	cfg := &config{
		workDir: t.TempDir(),
		envSpecMap: EnvSpecMap{
			"ENV_A": &EnvSpec{Name: "ENV_A"},
			"ENV_B": &EnvSpec{Name: "ENV_B"},
		},
	}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	// Test envs function - returns all defined env vars
	result, err := engine.Eval(`envs()`)
	assert.NoError(t, err)

	envMap, ok := result.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, "value_a", envMap["ENV_A"])
	assert.Equal(t, "value_b", envMap["ENV_B"])
}

func TestCelEngine_EnvWithoutEnvSpecMap(t *testing.T) {
	// When envSpecMap is nil, env() should work without validation (for patch_envs processing)
	assert.NoError(t, os.Setenv("TEST_VAR_NO_SPEC", "value123"))
	env.Reload()
	defer func() {
		assert.NoError(t, os.Unsetenv("TEST_VAR_NO_SPEC"))
	}()

	cfg := &config{
		workDir:    t.TempDir(),
		envSpecMap: nil, // No envSpecMap - env() calls not validated
	}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	result, err := engine.Eval(`env("TEST_VAR_NO_SPEC")`)
	assert.NoError(t, err)
	assert.Equal(t, "value123", result)
}

func TestCelEngine_EmbedFunction(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testContent := "Hello, CEL!"
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	assert.NoError(t, err)

	cfg := &config{workDir: tmpDir}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	// Test embed function
	result, err := engine.Eval(`embed("test.txt")`)
	assert.NoError(t, err)
	// Result should be base64 encoded
	assert.NotEmpty(t, result)

	// Test embed with empty string
	result, err = engine.Eval(`embed("")`)
	assert.NoError(t, err)
	assert.Equal(t, "", result)

	// Test embed with path traversal (should fail silently and return empty)
	result, err = engine.Eval(`embed("../../../etc/passwd")`)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestCelEngine_SecurityNoSideEffects(t *testing.T) {
	cfg := &config{workDir: t.TempDir()}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	// CEL should not allow arbitrary code execution
	// These expressions should fail to compile or evaluate

	dangerousExprs := []string{
		// No function calls to unknown functions
		`os.Exit(1)`,
		`exec("rm -rf /")`,
		// No variable assignment
		`x = 1`,
	}

	for _, expr := range dangerousExprs {
		t.Run(expr, func(t *testing.T) {
			_, err := engine.Eval(expr)
			// Should error on compile or evalExpr
			assert.Error(t, err, "expected error for dangerous expression: %s", expr)
		})
	}
}

func TestCelEngine_Ternary(t *testing.T) {
	cfg := &config{workDir: t.TempDir()}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	// Test ternary operator
	result, err := engine.Eval(`true ? "yes" : "no"`)
	assert.NoError(t, err)
	assert.Equal(t, "yes", result)

	result, err = engine.Eval(`false ? "yes" : "no"`)
	assert.NoError(t, err)
	assert.Equal(t, "no", result)
}

func TestCelEngine_StringOperations(t *testing.T) {
	cfg := &config{workDir: t.TempDir()}
	engine, err := newCelEngine(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		expr     string
		expected any
	}{
		{
			name:     "string contains",
			expr:     `"hello world".contains("world")`,
			expected: true,
		},
		{
			name:     "string startsWith",
			expr:     `"hello".startsWith("he")`,
			expected: true,
		},
		{
			name:     "string endsWith",
			expr:     `"hello".endsWith("lo")`,
			expected: true,
		},
		{
			name:     "string size",
			expr:     `size("hello")`,
			expected: int64(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Eval(tt.expr)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateCelFunction_AllSignatures(t *testing.T) {
	tests := []struct {
		name       string
		fn         any
		expr       string
		expected   any
		wantErr    bool
		expectNull bool // for fn() error success case
	}{
		{
			name:     "fn() T",
			fn:       func() string { return "hello" },
			expr:     `test_fn()`,
			expected: "hello",
		},
		{
			name:     "fn() int",
			fn:       func() int { return 42 },
			expr:     `test_fn()`,
			expected: int64(42),
		},
		{
			name: "fn() (T, error) - success",
			fn: func() (string, error) {
				return "success", nil
			},
			expr:     `test_fn()`,
			expected: "success",
		},
		{
			name: "fn() (T, error) - error",
			fn: func() (string, error) {
				return "", fmt.Errorf("test error")
			},
			expr:    `test_fn()`,
			wantErr: true,
		},
		{
			name: "fn() error - success",
			fn: func() error {
				return nil
			},
			expr:       `test_fn()`,
			expectNull: true,
		},
		{
			name: "fn() error - error",
			fn: func() error {
				return fmt.Errorf("test error")
			},
			expr:    `test_fn()`,
			wantErr: true,
		},
		{
			name:     "fn(T) R",
			fn:       func(s string) string { return "got:" + s },
			expr:     `test_fn("input")`,
			expected: "got:input",
		},
		{
			name:     "fn(int) int",
			fn:       func(n int64) int64 { return n * 2 },
			expr:     `test_fn(21)`,
			expected: int64(42),
		},
		{
			name: "fn(T) (R, error) - success",
			fn: func(s string) (string, error) {
				return "processed:" + s, nil
			},
			expr:     `test_fn("data")`,
			expected: "processed:data",
		},
		{
			name: "fn(T) (R, error) - error",
			fn: func(s string) (string, error) {
				return "", fmt.Errorf("processing failed")
			},
			expr:    `test_fn("data")`,
			wantErr: true,
		},
		{
			name: "fn(T) error - success",
			fn: func(s string) error {
				return nil
			},
			expr:       `test_fn("test")`,
			expectNull: true,
		},
		{
			name: "fn(T) error - error",
			fn: func(s string) error {
				return fmt.Errorf("validation failed: %s", s)
			},
			expr:    `test_fn("bad")`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh manager for each test
			mgr := NewManager()
			err := mgr.RegisterExprFunc("test_fn", tt.fn)
			assert.NoError(t, err)

			// Temporarily swap global manager
			oldMgr := globalManager
			globalManager = mgr
			defer func() { globalManager = oldMgr }()

			cfg := &config{workDir: t.TempDir()}
			engine, err := newCelEngine(cfg)
			assert.NoError(t, err)

			result, err := engine.Eval(tt.expr)
			if tt.wantErr {
				assert.Error(t, err)
			} else if tt.expectNull {
				assert.NoError(t, err)
				// NullValue returns structpb.NullValue(0)
				assert.NotNil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
