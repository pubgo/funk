package metaflags

import (
	"fmt"
	"strings"
	"time"
)

type StringValue struct {
	p    string
	name string
}

func (s *StringValue) Key() string { return s.name }

func (s *StringValue) get() any { return s.p }
func (s *StringValue) set(val any) error {
	switch v := val.(type) {
	case string:
		s.p = v
	case fmt.Stringer:
		s.p = v.String()
	default:
		s.p = fmt.Sprintf("%v", v)
	}
	return nil
}
func (s *StringValue) Get() string { return s.p }
func (s *StringValue) Set(val string) error {
	s.p = val
	return nil
}
func (s *StringValue) String() string { return s.p }
func String(name, value, usage string) *StringValue {
	s := &StringValue{p: value, name: name}
	defaultSet.AddFunc(name, s.get, s.set, usage)
	return s
}

type IntValue struct {
	p    *int64
	name string
}

func (i *IntValue) Key() string      { return i.name }
func (i *IntValue) Get() interface{} { return *i.p }
func (i *IntValue) Set(val interface{}) error {
	var x int64
	switch v := val.(type) {
	case int:
		x = int64(v)
	case int64:
		x = v
	case float64:
		x = int64(v)
	case string:
		fmt.Sscanf(v, "%d", &x)
	default:
		x = int64(fmt.Sprintf("%v", v)[0])
	}
	*i.p = x
	return nil
}
func (i *IntValue) String() string { return fmt.Sprintf("%d", *i.p) }
func Int(name string, value int64, usage string) Value {
	p := new(int64)
	*p = value
	return defaultSet.Add(name, &IntValue{p: p, name: name}, usage)
}

type FloatValue struct {
	p    *float64
	name string
}

func (f *FloatValue) Key() string      { return f.name }
func (f *FloatValue) Get() interface{} { return *f.p }
func (f *FloatValue) Set(val interface{}) error {
	var x float64
	switch v := val.(type) {
	case float64:
		x = v
	case float32:
		x = float64(v)
	case int:
		x = float64(v)
	case int64:
		x = float64(v)
	case string:
		fmt.Sscanf(v, "%f", &x)
	default:
		x = 0.0
	}
	*f.p = x
	return nil
}
func (f *FloatValue) String() string { return fmt.Sprintf("%g", *f.p) }
func Float(name string, value float64, usage string) Value {
	p := new(float64)
	*p = value
	return defaultSet.Add(name, &FloatValue{p: p, name: name}, usage)
}

type BoolValue struct {
	p    *bool
	name string
}

func (b *BoolValue) Key() string      { return b.name }
func (b *BoolValue) Get() interface{} { return *b.p }
func (b *BoolValue) Set(val interface{}) error {
	var x bool
	switch v := val.(type) {
	case bool:
		x = v
	case string:
		switch strings.ToLower(v) {
		case "true", "1", "on", "yes":
			x = true
		case "false", "0", "off", "no":
			x = false
		default:
			x = len(v) > 0
		}
	default:
		x = true
	}
	*b.p = x
	return nil
}
func (b *BoolValue) String() string { return fmt.Sprintf("%t", *b.p) }
func Bool(name string, value bool, usage string) Value {
	p := new(bool)
	*p = value
	return defaultSet.Add(name, &BoolValue{p: p, name: name}, usage)
}

type DurationValue struct {
	p    *time.Duration
	name string
}

func (d *DurationValue) Key() string      { return d.name }
func (d *DurationValue) Get() interface{} { return *d.p }
func (d *DurationValue) Set(val interface{}) error {
	var dur time.Duration
	switch v := val.(type) {
	case time.Duration:
		dur = v
	case string:
		var err error
		dur, err = time.ParseDuration(v)
		if err != nil {
			return err
		}
	case int:
		dur = time.Duration(v) * time.Second
	case float64:
		dur = time.Duration(v * float64(time.Second))
	default:
		dur = 0
	}
	*d.p = dur
	return nil
}
func (d *DurationValue) String() string { return (*d.p).String() }
func Duration(name string, value time.Duration, usage string) Value {
	p := new(time.Duration)
	*p = value
	return defaultSet.Add(name, &DurationValue{p: p, name: name}, usage)
}

type TimeValue struct {
	p    *time.Time
	name string
}

func (t *TimeValue) Key() string      { return t.name }
func (t *TimeValue) Get() interface{} { return *t.p }
func (t *TimeValue) Set(val interface{}) error {
	var tm time.Time
	switch v := val.(type) {
	case time.Time:
		tm = v
	case string:
		var err error
		tm, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
	case int64:
		tm = time.Unix(v, 0)
	default:
		return fmt.Errorf("cannot convert %T to time.Time", val)
	}
	*t.p = tm
	return nil
}
func (t *TimeValue) String() string { return (*t.p).Format(time.RFC3339) }
func Time(name string, value time.Time, usage string) Value {
	p := new(time.Time)
	*p = value
	return defaultSet.Add(name, &TimeValue{p: p, name: name}, usage)
}

type StringSliceValue struct {
	p    *[]string
	name string
}

func (s *StringSliceValue) Key() string      { return s.name }
func (s *StringSliceValue) Get() interface{} { return *s.p }
func (s *StringSliceValue) Set(val interface{}) error {
	var strs []string
	switch v := val.(type) {
	case []string:
		strs = v
	case []interface{}:
		strs = make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
	case string:
		parts := strings.Split(v, ",")
		for _, part := range parts {
			strs = append(strs, strings.TrimSpace(part))
		}
	default:
		strs = []string{fmt.Sprintf("%v", val)}
	}
	*s.p = strs
	return nil
}
func (s *StringSliceValue) String() string { return fmt.Sprintf("%v", *s.p) }
func StringSlice(name string, value []string, usage string) Value {
	p := new([]string)
	*p = append([]string(nil), value...)
	return defaultSet.Add(name, &StringSliceValue{p: p, name: name}, usage)
}

type StringMapValue struct {
	p    *map[string]string
	name string
}

func (m *StringMapValue) Key() string      { return m.name }
func (m *StringMapValue) Get() interface{} { return *m.p }
func (m *StringMapValue) Set(val interface{}) error {
	var mp map[string]string
	switch v := val.(type) {
	case map[string]string:
		mp = v
	case map[string]interface{}:
		mp = make(map[string]string)
		for k, vv := range v {
			mp[k] = fmt.Sprintf("%v", vv)
		}
	case string:
		mp = make(map[string]string)
		pairs := strings.Split(v, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				mp[kv[0]] = kv[1]
			}
		}
	default:
		mp = map[string]string{"value": fmt.Sprintf("%v", val)}
	}
	*m.p = mp
	return nil
}
func (m *StringMapValue) String() string { return fmt.Sprintf("%v", *m.p) }
func StringMap(name string, value map[string]string, usage string) Value {
	p := new(map[string]string)
	*p = make(map[string]string)
	for k, v := range value {
		(*p)[k] = v
	}
	return defaultSet.Add(name, &StringMapValue{p: p, name: name}, usage)
}

type FuncValue struct {
	get  func() any
	set  func(val any) error
	name string
}

func (f *FuncValue) Key() string       { return f.name }
func (f *FuncValue) Get() any          { return f.get() }
func (f *FuncValue) Set(val any) error { return f.set(val) }
func (f *FuncValue) String() string    { return fmt.Sprintf("%v", f.get()) }

func Func(name string, getter func() any, setter func(val any) error, usage string) Value {
	return defaultSet.Add(name, &FuncValue{get: getter, set: setter, name: name}, usage)
}

func Any(name string, value any, usage string) Value {
	anyVal := value
	return Func(
		name,
		func() any { return anyVal }, func(val any) error { anyVal = val; return nil },
		usage,
	)
}
