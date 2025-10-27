package features

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k0kubun/pp/v3"
)

var _ Value = (*baseValue[any])(nil)

func newBase[T any](f *Feature, name string, value T, usage string, typ ValueType, tags []map[string]any, set func(s string) (T, error), getString func(val T) string) *baseValue[T] {
	base := &baseValue[T]{val: value, set: set, getString: getString, typ: typ}
	base.ff = f.AddFunc(name, usage, base, tags...)
	return base
}

type baseValue[T any] struct {
	ff        *Flag
	typ       ValueType
	val       T
	set       func(s string) (T, error)
	getString func(val T) string
}

func (b *baseValue[T]) Name() string    { return b.ff.Name }
func (b *baseValue[T]) Type() ValueType { return b.typ }
func (b *baseValue[T]) Get() any        { return b.val }
func (b *baseValue[T]) GetValue() T     { return b.val }
func (b *baseValue[T]) Set(s string) error {
	val, err := b.set(s)
	if err != nil {
		return fmt.Errorf("faield to set value, value=%s err=%w", s, err)
	}

	b.val = val
	return nil
}

func (b *baseValue[T]) String() string {
	if b.getString == nil {
		return fmt.Sprintf("%v", b.val)
	}

	return b.getString(b.val)
}

type StringValue struct {
	*baseValue[string]
}

func String(name, value, usage string, tags ...map[string]any) StringValue {
	base := newBase(
		defaultFeature,
		name,
		value,
		usage,
		StringType,
		tags,
		func(s string) (string, error) { return s, nil },
		func(val string) string { return val },
	)
	return StringValue{baseValue: base}
}

type IntValue struct {
	*baseValue[int64]
}

func Int(name string, value int64, usage string, tags ...map[string]any) IntValue {
	base := newBase(
		defaultFeature,
		name,
		value,
		usage,
		IntType,
		tags,
		func(s string) (val int64, err error) {
			_, err = fmt.Sscanf(s, "%d", &val)
			return val, err
		},
		nil,
	)
	return IntValue{baseValue: base}
}

type FloatValue struct {
	*baseValue[float64]
}

func Float(name string, value float64, usage string, tags ...map[string]any) FloatValue {
	base := newBase(
		defaultFeature,
		name,
		value,
		usage,
		FloatType,
		tags,
		func(s string) (val float64, err error) {
			_, err = fmt.Sscanf(s, "%f", &val)
			return val, err
		},
		func(val float64) string { return fmt.Sprintf("%f", val) },
	)
	return FloatValue{baseValue: base}
}

type BoolValue struct {
	*baseValue[bool]
}

func Bool(name string, value bool, usage string, tags ...map[string]any) BoolValue {
	base := newBase(
		defaultFeature,
		name,
		value,
		usage,
		BoolType,
		tags,
		func(s string) (val bool, err error) {
			switch strings.ToLower(s) {
			case "true", "1", "on", "yes":
				return true, nil
			case "false", "0", "off", "no":
				return false, nil
			default:
				return len(s) > 0, nil
			}
		},
		func(val bool) string { return fmt.Sprintf("%v", val) },
	)
	return BoolValue{baseValue: base}
}

type JsonValue[T any] struct {
	*baseValue[T]
}

func Json[T any](name string, value T, usage string, tags ...map[string]any) JsonValue[T] {
	base := newBase[T](
		defaultFeature,
		name,
		value,
		usage,
		JsonType,
		tags,
		func(s string) (val T, err error) {
			return val, json.Unmarshal([]byte(s), &val)
		},
		func(val T) string {
			data, err := json.Marshal(val)
			if err != nil {
				_, _ = pp.Println(val)
				return err.Error()
			}
			return string(data)
		},
	)
	return JsonValue[T]{baseValue: base}
}
