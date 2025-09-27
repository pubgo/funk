package monster

import (
	"sync"
)

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

// Monster manages all registered values
type Monster struct {
	mutex sync.RWMutex
	m     map[string]*Entry
}

var defaultMonster = NewMonster()

// NewMonster creates a new Monster instance
func NewMonster() *Monster {
	return &Monster{
		m: make(map[string]*Entry),
	}
}

// AddFunc registers a new value with getter, setter, usage, and optional tags
func (m *Monster) AddFunc(name string, get Getter, set Setter, usage string, tags ...map[string]any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
func (m *Monster) Lookup(name string) *Entry {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.m[name]
}

// VisitAll calls fn for each entry
func (m *Monster) VisitAll(fn func(*Entry)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, e := range m.m {
		fn(e)
	}
}

func Register(name string, get Getter, set Setter, usage string, tags ...map[string]any) {
	defaultMonster.AddFunc(name, get, set, usage, tags...)
}

func Lookup(name string) *Entry {
	return defaultMonster.Lookup(name)
}

func VisitAll(fn func(*Entry)) { defaultMonster.VisitAll(fn) }
