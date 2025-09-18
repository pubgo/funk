package monitor

import (
	"fmt"
	"strings"
	"time"
)

// mergeTags safely copies optional tags
func mergeTags(maps ...map[string]any) map[string]any {
	if len(maps) == 0 || maps[0] == nil {
		return make(map[string]any)
	}
	m := make(map[string]any)
	for k, v := range maps[0] {
		m[k] = v
	}
	return m
}

type StringValue struct {
	p    string
	name string
}

func (s *StringValue) Key() string { return s.name }
func (s *StringValue) get() any    { return s.p }
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
func (s *StringValue) Get() string          { return s.p }
func (s *StringValue) Set(val string) error { return s.set(val) }
func (s *StringValue) String() string       { return s.p }

func String(name, value, usage string, tags ...map[string]any) *StringValue {
	s := &StringValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, s.get, s.set, usage, tagCopy)
	return s
}

type IntValue struct {
	p    int64
	name string
}

func (i *IntValue) Key() string { return i.name }
func (i *IntValue) get() any    { return i.p }
func (i *IntValue) set(val any) error {
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
		x = 0
	}
	i.p = x
	return nil
}
func (i *IntValue) Get() int64          { return i.p }
func (i *IntValue) Set(val int64) error { return i.set(val) }
func (i *IntValue) String() string      { return fmt.Sprintf("%d", i.p) }

func Int(name string, value int64, usage string, tags ...map[string]any) *IntValue {
	i := &IntValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, i.get, i.set, usage, tagCopy)
	return i
}

type FloatValue struct {
	p    float64
	name string
}

func (f *FloatValue) Key() string { return f.name }
func (f *FloatValue) get() any    { return f.p }
func (f *FloatValue) set(val any) error {
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
	f.p = x
	return nil
}
func (f *FloatValue) Get() float64          { return f.p }
func (f *FloatValue) Set(val float64) error { return f.set(val) }
func (f *FloatValue) String() string        { return fmt.Sprintf("%g", f.p) }

func Float(name string, value float64, usage string, tags ...map[string]any) *FloatValue {
	f := &FloatValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, f.get, f.set, usage, tagCopy)
	return f
}

type BoolValue struct {
	p    bool
	name string
}

func (b *BoolValue) Key() string { return b.name }
func (b *BoolValue) get() any    { return b.p }
func (b *BoolValue) set(val any) error {
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
	b.p = x
	return nil
}
func (b *BoolValue) Get() bool          { return b.p }
func (b *BoolValue) Set(val bool) error { return b.set(val) }
func (b *BoolValue) String() string     { return fmt.Sprintf("%t", b.p) }

func Bool(name string, value bool, usage string, tags ...map[string]any) *BoolValue {
	b := &BoolValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, b.get, b.set, usage, tagCopy)
	return b
}

type DurationValue struct {
	p    time.Duration
	name string
}

func (d *DurationValue) Key() string { return d.name }
func (d *DurationValue) get() any    { return d.p }
func (d *DurationValue) set(val any) error {
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
	default:
		return fmt.Errorf("cannot convert %T to time.Duration", val)
	}
	d.p = dur
	return nil
}
func (d *DurationValue) Get() time.Duration          { return d.p }
func (d *DurationValue) Set(val time.Duration) error { return d.set(val) }
func (d *DurationValue) String() string              { return d.p.String() }

func Duration(name string, value time.Duration, usage string, tags ...map[string]any) *DurationValue {
	d := &DurationValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, d.get, d.set, usage, tagCopy)
	return d
}

type TimeValue struct {
	p    time.Time
	name string
}

func (t *TimeValue) Key() string { return t.name }
func (t *TimeValue) get() any    { return t.p }
func (t *TimeValue) set(val any) error {
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
	t.p = tm
	return nil
}
func (t *TimeValue) Get() time.Time          { return t.p }
func (t *TimeValue) Set(val time.Time) error { return t.set(val) }
func (t *TimeValue) String() string          { return t.p.Format(time.RFC3339) }

func Time(name string, value time.Time, usage string, tags ...map[string]any) *TimeValue {
	t := &TimeValue{p: value, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, t.get, t.set, usage, tagCopy)
	return t
}

type StringSliceValue struct {
	p    []string
	name string
}

func (s *StringSliceValue) Key() string { return s.name }
func (s *StringSliceValue) get() any    { return s.p }
func (s *StringSliceValue) set(val any) error {
	var strs []string
	switch v := val.(type) {
	case []string:
		strs = v
	case []interface{}:
		strs = make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
	default:
		strs = []string{fmt.Sprintf("%v", val)}
	}
	s.p = strs
	return nil
}
func (s *StringSliceValue) Get() []string          { return s.p }
func (s *StringSliceValue) Set(val []string) error { return s.set(val) }
func (s *StringSliceValue) String() string         { return fmt.Sprintf("%v", s.p) }

func StringSlice(name string, value []string, usage string, tags ...map[string]any) *StringSliceValue {
	cp := append([]string(nil), value...)
	s := &StringSliceValue{p: cp, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, s.get, s.set, usage, tagCopy)
	return s
}

type StringMapValue struct {
	p    map[string]string
	name string
}

func (m *StringMapValue) Key() string { return m.name }
func (m *StringMapValue) get() any    { return m.p }
func (m *StringMapValue) set(val any) error {
	var mp map[string]string
	switch v := val.(type) {
	case map[string]string:
		mp = v
	case map[string]interface{}:
		mp = make(map[string]string)
		for k, vv := range v {
			mp[k] = fmt.Sprintf("%v", vv)
		}
	default:
		return fmt.Errorf("cannot convert %T to map[string]string", val)
	}
	m.p = mp
	return nil
}
func (m *StringMapValue) Get() map[string]string          { return m.p }
func (m *StringMapValue) Set(val map[string]string) error { return m.set(val) }
func (m *StringMapValue) String() string                  { return fmt.Sprintf("%v", m.p) }

func StringMap(name string, value map[string]string, usage string, tags ...map[string]any) *StringMapValue {
	cp := make(map[string]string)
	for k, v := range value {
		cp[k] = v
	}
	m := &StringMapValue{p: cp, name: name}
	tagCopy := mergeTags(tags...)
	defaultMonitor.AddFunc(name, m.get, m.set, usage, tagCopy)
	return m
}
