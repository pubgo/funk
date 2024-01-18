package config

import (
	"github.com/goccy/go-json"
	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = (*Node)(nil)
var _ yaml.Marshaler = (*Node)(nil)
var _ json.Marshaler = (*Node)(nil)

type Node struct {
	maps  map[string]any
	value *yaml.Node
}

func (c *Node) YamlNode() *yaml.Node {
	return c.value
}

func (c *Node) MarshalYAML() (interface{}, error) {
	return c.value, nil
}

func (c *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.maps)
}

func (c *Node) UnmarshalYAML(value *yaml.Node) error {
	if c.maps == nil {
		c.maps = make(map[string]any)
	}

	if err := value.Decode(&c.maps); err != nil {
		return err
	}

	c.value = value

	return nil
}

func (c *Node) IsNil() bool {
	return c.value == nil
}

func (c *Node) Get(key string) any {
	if c.maps == nil {
		return nil
	}

	return c.maps[key]
}

func (c *Node) Decode(val any) error {
	if c.IsNil() {
		return nil
	}

	return c.value.Decode(val)
}

type ListOrMap[T any] []T

// MarshalYAML implements the yaml.Marshaler interface.
func (ts *ListOrMap[T]) MarshalYAML() (any, error) {
	if len(*ts) == 1 {
		return (*ts)[0], nil
	} else {
		return []T(*ts), nil
	}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (ts *ListOrMap[T]) UnmarshalYAML(value *yaml.Node) error {
	return unmarshalOneOrList((*[]T)(ts), value)
}
