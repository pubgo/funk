package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurePath(t *testing.T) {
	baseDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(baseDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// Create a subdirectory with file
	subDir := filepath.Join(baseDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)
	subFile := filepath.Join(subDir, "sub.txt")
	err = os.WriteFile(subFile, []byte("sub"), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		base      string
		target    string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid file in base",
			base:    baseDir,
			target:  "test.txt",
			wantErr: false,
		},
		{
			name:    "valid file in subdir",
			base:    baseDir,
			target:  "subdir/sub.txt",
			wantErr: false,
		},
		{
			name:      "path traversal attack",
			base:      baseDir,
			target:    "../../../etc/passwd",
			wantErr:   true,
			errSubstr: "path traversal detected",
		},
		{
			name:      "path traversal with subdir",
			base:      baseDir,
			target:    "subdir/../../etc/passwd",
			wantErr:   true,
			errSubstr: "path traversal detected",
		},
		{
			name:      "empty name",
			base:      baseDir,
			target:    "",
			wantErr:   true,
			errSubstr: "file name is empty",
		},
		{
			name:      "absolute path escape",
			base:      baseDir,
			target:    "/etc/passwd",
			wantErr:   true,
			errSubstr: "absolute path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := securePath(tt.base, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestMaskSensitiveData(t *testing.T) {
	// Short data should not be truncated
	short := "password: secret123"
	result := maskSensitiveData(short)
	assert.Equal(t, short, result)

	// Long data should be truncated
	long := make([]byte, 3000)
	for i := range long {
		long[i] = 'a'
	}
	result = maskSensitiveData(string(long))
	assert.Contains(t, result, "[truncated for safety]")
	assert.Less(t, len(result), len(long))
}

func TestIsSensitiveKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"password", true},
		{"db_password", true},
		{"PASSWORD", true},
		{"api_key", true},
		{"apiKey", true},
		{"secret", true},
		{"token", true},
		{"auth_token", true},
		{"dsn", true},
		{"username", false},
		{"host", false},
		{"port", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := isSensitiveKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvSpecValidation(t *testing.T) {
	tests := []struct {
		name    string
		spec    EnvSpec
		wantErr bool
	}{
		{
			name: "required missing",
			spec: EnvSpec{
				Name: "TEST_REQUIRED",
				Rule: "required",
			},
			wantErr: true,
		},
		{
			name: "required present",
			spec: EnvSpec{
				Name:    "TEST_REQUIRED_OK",
				Rule:    "required",
				Default: "value",
			},
			wantErr: false,
		},
		{
			name: "numeric valid",
			spec: EnvSpec{
				Name:    "TEST_NUMERIC",
				Rule:    "numeric",
				Default: "123",
			},
			wantErr: false,
		},
		{
			name: "numeric invalid",
			spec: EnvSpec{
				Name:    "TEST_NUMERIC_BAD",
				Rule:    "numeric",
				Default: "abc",
			},
			wantErr: true,
		},
		{
			name: "email valid",
			spec: EnvSpec{
				Name:    "TEST_EMAIL",
				Rule:    "email",
				Default: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "email invalid",
			spec: EnvSpec{
				Name:    "TEST_EMAIL_BAD",
				Rule:    "email",
				Default: "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "url valid",
			spec: EnvSpec{
				Name:    "TEST_URL",
				Rule:    "url",
				Default: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "url invalid",
			spec: EnvSpec{
				Name:    "TEST_URL_BAD",
				Rule:    "url",
				Default: "not-a-url",
			},
			wantErr: true,
		},
		{
			name: "uuid valid",
			spec: EnvSpec{
				Name:    "TEST_UUID",
				Rule:    "uuid",
				Default: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "uuid invalid",
			spec: EnvSpec{
				Name:    "TEST_UUID_BAD",
				Rule:    "uuid",
				Default: "not-a-uuid",
			},
			wantErr: true,
		},
		{
			name: "ip valid",
			spec: EnvSpec{
				Name:    "TEST_IP",
				Rule:    "ip",
				Default: "192.168.1.1",
			},
			wantErr: false,
		},
		{
			name: "min,max valid",
			spec: EnvSpec{
				Name:    "TEST_LEN",
				Rule:    "min=3,max=10",
				Default: "hello",
			},
			wantErr: false,
		},
		{
			name: "min too short",
			spec: EnvSpec{
				Name:    "TEST_LEN_SHORT",
				Rule:    "min=5",
				Default: "hi",
			},
			wantErr: true,
		},
		{
			name: "no rule - always pass",
			spec: EnvSpec{
				Name:    "TEST_NO_RULE",
				Default: "anything",
			},
			wantErr: false,
		},
		{
			name: "boolean valid",
			spec: EnvSpec{
				Name:    "TEST_BOOL",
				Rule:    "boolean",
				Default: "true",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManagerThreadSafety(t *testing.T) {
	m := &Manager{
		exprFns: make(map[string]any),
	}

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			m.SetPath("/path/" + string(rune('a'+n)))
			_ = m.GetPath()
			m.SetDir("/dir/" + string(rune('a'+n)))
			_ = m.GetDir()
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRegisterExpr(t *testing.T) {
	// Reset for test
	oldManager := globalManager
	globalManager = &Manager{
		exprFns: make(map[string]any),
	}
	defer func() { globalManager = oldManager }()

	// First registration should succeed
	err := RegisterExpr("myFunc", func() string { return "test" })
	assert.NoError(t, err)

	// Duplicate registration should fail
	err = RegisterExpr("myFunc", func() string { return "test2" })
	assert.Error(t, err)
	assert.IsType(t, &ExprExistsError{}, err)
}
