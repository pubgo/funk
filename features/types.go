package features

import (
	"encoding/json"
	"fmt"
	"strings"
)

type StringValue struct {
	val string
	ff  *Flag
}

func (s *StringValue) Type() string         { return "string" }
func (s *StringValue) Name() string         { return s.ff.Name }
func (s *StringValue) Get() any             { return s.val }
func (s *StringValue) GetValue() string     { return s.val }
func (s *StringValue) Set(val string) error { s.val = val; return nil }
func (s *StringValue) String() string       { return s.val }

func String(name, value, usage string, tags ...map[string]any) *StringValue {
	s := &StringValue{val: value}
	s.ff = defaultFeature.AddFunc(name, usage, s, tags...)
	return s
}

type IntValue struct {
	val int64
	ff  *Flag
}

func (i *IntValue) Type() string    { return "int" }
func (i *IntValue) Name() string    { return i.ff.Name }
func (i *IntValue) Get() any        { return i.val }
func (i *IntValue) GetValue() int64 { return i.val }
func (i *IntValue) String() string  { return fmt.Sprintf("%d", i.val) }
func (i *IntValue) Set(val string) error {
	_, err := fmt.Sscanf(val, "%d", &i.val)
	return err
}

func Int(name string, value int64, usage string, tags ...map[string]any) *IntValue {
	i := &IntValue{val: value}
	i.ff = defaultFeature.AddFunc(name, usage, i, tags...)
	return i
}

type FloatValue struct {
	val float64
	ff  *Flag
}

func (f *FloatValue) Type() string   { return "float" }
func (f *FloatValue) Name() string   { return f.ff.Name }
func (f *FloatValue) Get() any       { return f.val }
func (f *FloatValue) String() string { return fmt.Sprintf("%f", f.val) }
func (f *FloatValue) Set(val string) error {
	_, err := fmt.Sscanf(val, "%f", &f.val)
	return err
}

func Float(name string, value float64, usage string, tags ...map[string]any) *FloatValue {
	f := &FloatValue{val: value}
	f.ff = defaultFeature.AddFunc(name, usage, f, tags...)
	return f
}

type BoolValue struct {
	val bool
	ff  *Flag
}

func (b *BoolValue) Type() string   { return "bool" }
func (b *BoolValue) Name() string   { return b.ff.Name }
func (b *BoolValue) Get() any       { return b.val }
func (b *BoolValue) GetValue() bool { return b.val }
func (b *BoolValue) String() string { return fmt.Sprintf("%v", b.val) }
func (b *BoolValue) Set(val string) error {
	switch strings.ToLower(val) {
	case "true", "1", "on", "yes":
		b.val = true
	case "false", "0", "off", "no":
		b.val = false
	default:
		b.val = len(val) > 0
	}
	return nil
}

func Bool(name string, value bool, usage string, tags ...map[string]any) *BoolValue {
	b := &BoolValue{val: value}
	b.ff = defaultFeature.AddFunc(name, usage, b, tags...)
	return b
}

type JsonValue[T any] struct {
	val T
	ff  *Flag
}

func (t *JsonValue[T]) Type() string { return "json" }
func (t *JsonValue[T]) Name() string { return t.ff.Name }
func (t *JsonValue[T]) GetValue() T  { return t.val }
func (t *JsonValue[T]) Get() any     { return t.val }
func (t *JsonValue[T]) String() string {
	data, err := json.Marshal(t.val)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
func (t *JsonValue[T]) Set(val string) error {
	return json.Unmarshal([]byte(val), &t.val)
}
func Json[T any](name string, value T, usage string, tags ...map[string]any) *JsonValue[T] {
	t := &JsonValue[T]{val: value}
	t.ff = defaultFeature.AddFunc(name, usage, t, tags...)
	return t
}
