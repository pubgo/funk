package typex

import (
	yaml "gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/errors"
)

type YamlListType[T any] []T

func (p *YamlListType[T]) UnmarshalYAML(value *yaml.Node) error {
	if value.IsZero() {
		return nil
	}

	switch value.Kind {
	case yaml.ScalarNode, yaml.MappingNode:
		var data T
		if err := value.Decode(&data); err != nil {
			return errors.WrapCaller(err)
		}
		*p = []T{data}
		return nil
	case yaml.SequenceNode:
		var data []T
		if err := value.Decode(&data); err != nil {
			return errors.WrapCaller(err)
		}
		*p = data
		return nil
	default:
		var val any
		assert.Exit(value.Decode(&val))
		return errors.Errorf("yaml kind type error, kind=%v data=%v", value.Kind, val)
	}
}
