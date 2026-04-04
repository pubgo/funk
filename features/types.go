package features

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Ensure baseValue implements Value interface
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

func (b *baseValue[T]) Name() string { return b.ff.Name }
func (b *baseValue[T]) Type() string { return b.typ.String() }
func (b *baseValue[T]) Value() any   { return b.val }
func (b *baseValue[T]) Set(s string) error {
	val, err := b.set(s)
	if err != nil {
		return fmt.Errorf("failed to set value, value=%s err=%w", s, err)
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

func (v StringValue) Value() string { return v.val }

func String(name, value, usage string, tags ...map[string]any) StringValue {
	return StringValue{baseValue: newBase(
		defaultFeature,
		name,
		value,
		usage,
		StringType,
		tags,
		func(s string) (string, error) { return s, nil },
		func(val string) string { return val },
	)}
}

type IntValue struct {
	*baseValue[int64]
}

func (v IntValue) Value() int64 { return v.val }

func Int(name string, value int64, usage string, tags ...map[string]any) IntValue {
	return IntValue{baseValue: newBase(
		defaultFeature,
		name,
		value,
		usage,
		IntType,
		tags,
		func(s string) (val int64, err error) {
			_, err = fmt.Sscanf(s, "%d", &val)
			if err != nil {
				return val, fmt.Errorf("failed to parse int value, str=%s err=%w", s, err)
			}
			return val, nil
		},
		func(val int64) string { return fmt.Sprintf("%d", val) },
	)}
}

type FloatValue struct {
	*baseValue[float64]
}

func (v FloatValue) Value() float64 { return v.val }

func Float(name string, value float64, usage string, tags ...map[string]any) FloatValue {
	return FloatValue{baseValue: newBase(
		defaultFeature,
		name,
		value,
		usage,
		FloatType,
		tags,
		func(s string) (val float64, err error) {
			_, err = fmt.Sscanf(s, "%f", &val)
			if err != nil {
				return val, fmt.Errorf("failed to parse float value, str=%s err=%w", s, err)
			}
			return val, nil
		},
		func(val float64) string { return fmt.Sprintf("%f", val) },
	)}
}

type BoolValue struct {
	*baseValue[bool]
}

func (v BoolValue) Value() bool { return v.val }

func Bool(name string, value bool, usage string, tags ...map[string]any) BoolValue {
	return BoolValue{baseValue: newBase(
		defaultFeature,
		name,
		value,
		usage,
		BoolType,
		tags,
		func(s string) (val bool, err error) {
			switch strings.ToLower(s) {
			case "true", "1", "on", "yes", "ok":
				return true, nil
			case "false", "0", "off", "no", "fail":
				return false, nil
			default:
				return len(s) > 0, nil
			}
		},
		func(val bool) string { return fmt.Sprintf("%v", val) },
	)}
}

type JsonValue[T any] struct {
	*baseValue[T]
}

func (v JsonValue[T]) Value() T { return v.val }

func Json[T any](name string, value T, usage string, tags ...map[string]any) JsonValue[T] {
	return JsonValue[T]{baseValue: newBase[T](
		defaultFeature,
		name,
		value,
		usage,
		JsonType,
		tags,
		func(s string) (val T, err error) {
			err = json.Unmarshal([]byte(s), &val)
			if err != nil {
				return val, fmt.Errorf("failed to unmarshal json, str=%s err=%w", s, err)
			}
			return val, nil
		},
		func(val T) string {
			data, err := json.Marshal(val)
			if err != nil {
				return fmt.Sprintf("failed to marshal json, val=%v err=%s", val, err.Error())
			}
			return string(data)
		},
	)}
}
