package features

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/spf13/pflag"
)

type ValueType string

func (vt ValueType) String() string {
	switch vt {
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
	pflag.Value
	Value() any
}

var _ json.Marshaler = (*Flag)(nil)

type Flag struct {
	Name       string         `json:"name"`
	Usage      string         `json:"usage"`
	Value      Value          `json:"-"`
	Deprecated bool           `json:"deprecated"`
	Tags       map[string]any `json:"tags"`
}

func (f Flag) MarshalJSON() ([]byte, error) {
	data := make(map[string]any)
	data["name"] = f.Name
	data["usage"] = f.Usage
	data["value"] = f.Value.Value()
	data["type"] = f.Value.Type()
	data["deprecated"] = f.Deprecated
	return json.Marshal(data)
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
		panic(fmt.Sprintf("feature flag already exists, name:%s", name))
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
	m := make(map[string]any)
	for _, mm := range maps {
		for k, v := range mm {
			m[k] = v
		}
	}
	return m
}
