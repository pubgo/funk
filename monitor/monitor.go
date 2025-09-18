package monitor

import "sync"

// Getter returns the current value as any
type Getter func() any

// Setter updates the value from any
type Setter func(any) error

// Entry holds metadata about a registered value
type Entry struct {
	Name   string
	Getter Getter
	Setter Setter
	Usage  string
	Tags   map[string]any // 自定义元数据标签
}

// Monitor manages all registered values
type Monitor struct {
	m  map[string]*Entry
	mu sync.RWMutex
}

var defaultMonitor = NewMonitor()

// NewMonitor creates a new Monitor instance
func NewMonitor() *Monitor {
	return &Monitor{
		m: make(map[string]*Entry),
	}
}

// AddFunc registers a new value with getter, setter, usage, and optional tags
func (m *Monitor) AddFunc(name string, get Getter, set Setter, usage string, tags ...map[string]any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tagMap := make(map[string]any)
	if len(tags) > 0 && tags[0] != nil {
		for k, v := range tags[0] {
			tagMap[k] = v
		}
	}

	m.m[name] = &Entry{
		Name:   name,
		Getter: get,
		Setter: set,
		Usage:  usage,
		Tags:   tagMap,
	}
}

// Lookup returns the entry by name
func (m *Monitor) Lookup(name string) *Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[name]
}

// VisitAll calls fn for each entry
func (m *Monitor) VisitAll(fn func(*Entry)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, e := range m.m {
		fn(e)
	}
}

func AddFunc(name string, get Getter, set Setter, usage string, tags ...map[string]any) {
	defaultMonitor.AddFunc(name, get, set, usage, tags...)
}

func Lookup(name string) *Entry {
	return defaultMonitor.Lookup(name)
}

func VisitAll(fn func(*Entry)) {
	defaultMonitor.VisitAll(fn)
}
