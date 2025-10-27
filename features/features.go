package features

import (
	"fmt"
	"sync"
)

type ValueType string

func (typ ValueType) String() string {
	switch typ {
	case StringType:
		return "string"
	case IntType:
		return "int"
	case FloatType:
		return "float"
	case BoolType:
		return "bool"
	case JsonType:
		return "json"
	default:
		return "unknown"
	}
}

const (
	StringType ValueType = "string"
	IntType    ValueType = "int"
	FloatType  ValueType = "float"
	BoolType   ValueType = "bool"
	JsonType   ValueType = "json"
)

type Value interface {
	String() string
	Type() ValueType
	Set(string) error
	Get() any
}

type Flag struct {
	Name       string
	Usage      string
	Value      Value
	Deprecated bool
	Tags       map[string]any
}

type Feature struct {
	mutex sync.RWMutex
	flags map[string]*Flag
}

// NewFeature creates a new Feature instance
func NewFeature() *Feature {
	return &Feature{
		flags: make(map[string]*Flag),
	}
}

// AddFunc registers a new value with getter, setter, usage, and optional tags
func (m *Feature) AddFunc(name, usage string, value Value, tags ...map[string]any) *Flag {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.flags[name] != nil {
		panic(fmt.Sprintf("flag already exists, name:%s", name))
	}

	ff := &Flag{Name: name, Usage: usage, Value: value, Tags: mergeTags(tags...)}
	m.flags[name] = ff
	return ff
}

// Lookup returns the entry by name
func (m *Feature) Lookup(name string) *Flag {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.flags[name]
}

// VisitAll calls fn for each entry
func (m *Feature) VisitAll(fn func(*Flag)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, e := range m.flags {
		fn(e)
	}
}

var defaultFeature = NewFeature()

func Register(name, usage string, value Value, tags ...map[string]any) {
	defaultFeature.AddFunc(name, usage, value, tags...)
}

func Lookup(name string) *Flag {
	return defaultFeature.Lookup(name)
}

func VisitAll(fn func(*Flag)) { defaultFeature.VisitAll(fn) }

// mergeTags safely copies optional tags
func mergeTags(maps ...map[string]any) map[string]any {
	if len(maps) == 0 || maps[0] == nil {
		return make(map[string]any)
	}

	m := make(map[string]any)
	for _, mm := range maps {
		for k, v := range mm {
			m[k] = v
		}
	}
	return m
}
