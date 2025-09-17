package metaflags

import (
	"github.com/pubgo/funk/assert"
	"sync"
)

// Value is the interface to the dynamic value stored in metadata.
type Value interface {
	Get() any
	Set(value any) error
	String() string
	Key() string // 返回该值的名称
}

// Entry holds metadata about a registered value
type Entry struct {
	Name  string
	Value Value
	Usage string
}

// FlagSet holds a set of metadata values
type FlagSet struct {
	m  map[string]*Entry
	mu sync.RWMutex
}

var defaultSet = NewFlagSet()

func NewFlagSet() *FlagSet { return &FlagSet{m: make(map[string]*Entry)} }

func (f *FlagSet) AddFunc(name string, getter func() any, setter func(val any) error, usage string) Value {
	return f.Add(name, &FuncValue{get: getter, set: setter, name: name}, usage)
}

// Add registers a new metadata entry and injects name if supported
func (f *FlagSet) Add(name string, value Value, usage string) Value {
	f.mu.Lock()
	defer f.mu.Unlock()

	assert.If(value == nil, "flag value is nil")
	assert.If(name == "", "flag name is empty")
	assert.If(f.m[name] != nil, "flag '%s' already registered", name)

	f.m[name] = &Entry{
		Name:  name,
		Value: value,
		Usage: usage,
	}
	return value
}

// Lookup returns the Entry for the given name
func (f *FlagSet) Lookup(name string) *Entry {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.m[name]
}

// VisitAll calls fn for each Entry
func (f *FlagSet) VisitAll(fn func(*Entry)) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, entry := range f.m {
		fn(entry)
	}
}

// Default accessors
func Lookup(name string) *Entry { return defaultSet.Lookup(name) }
func VisitAll(fn func(*Entry))  { defaultSet.VisitAll(fn) }
func Add(name string, val Value, u string) Value {
	return defaultSet.Add(name, val, u)
}
