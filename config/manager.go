package config

import (
	"sync"
)

// globalManager is the default config manager instance with thread-safe access
var globalManager = NewManager()

// NewManager creates a new Manager instance
func NewManager() *Manager {
	return &Manager{
		exprFns: make(map[string]any),
	}
}

// Manager provides thread-safe configuration management
type Manager struct {
	mu         sync.RWMutex
	configDir  string
	configPath string
	exprFns    map[string]any
	data       any
	envMap     EnvSpecMap
}

func (m *Manager) SetEnvMap(data EnvSpecMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.envMap = data
}

func (m *Manager) GetEnvMap() EnvSpecMap {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.envMap
}

func (m *Manager) SetConfigData(data any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = data
}

func (m *Manager) GetConfigData() any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data
}

// SetPath sets the config path in a thread-safe manner
func (m *Manager) SetPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configPath = path
}

// GetPath returns the config path in a thread-safe manner
func (m *Manager) GetPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.configPath
}

// SetDir sets the config directory in a thread-safe manner
func (m *Manager) SetDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configDir = dir
}

// GetDir returns the config directory in a thread-safe manner
func (m *Manager) GetDir() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.configDir
}

// RegisterExprFunc registers a custom expression function in a thread-safe manner
func (m *Manager) RegisterExprFunc(name string, fn any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.exprFns[name] != nil {
		return &ExprExistsError{Name: name}
	}
	m.exprFns[name] = fn
	return nil
}

// GetExprFns returns a copy of registered expression functions
func (m *Manager) GetExprFns() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]any, len(m.exprFns))
	for k, v := range m.exprFns {
		result[k] = v
	}
	return result
}

// ExprExistsError is returned when trying to register a duplicate expression
type ExprExistsError struct {
	Name string
}

func (e *ExprExistsError) Error() string {
	return "expr function already exists: " + e.Name
}

// Global convenience functions that use globalManager

// Global returns the global config manager instance
func Global() *Manager {
	return globalManager
}
