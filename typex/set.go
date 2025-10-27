package typex

import (
	"reflect"
	"sync"

	"go.uber.org/atomic"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
)

func SetOf(val ...any) *Set {
	s := &Set{}
	for i := range val {
		s.Add(val[i])
	}
	return s
}

type Set struct {
	m     sync.Map
	count atomic.Uint32
}

func (t *Set) Has(v any) bool { _, ok := t.m.Load(v); return ok }
func (t *Set) Len() uint32    { return t.count.Load() }

func (t *Set) Map(data any) (err error) {
	defer recovery.Err(&err)

	vd := reflect.ValueOf(data)
	assert.If(vd.Kind() != reflect.Ptr, "[data] should be ptr type")
	vd = vd.Elem()

	dt := reflect.MakeSlice(vd.Type(), 0, int(t.count.Load()))
	t.m.Range(func(key, _ any) bool {
		dt = reflect.AppendSlice(dt, reflect.ValueOf(key))
		return true
	})
	vd.Set(dt)

	return nil
}

func (t *Set) Add(v any) {
	_, ok := t.m.LoadOrStore(v, struct{}{})
	if !ok {
		t.count.Inc()
	}
}

func (t *Set) List() (val []any) {
	t.m.Range(func(key, _ any) bool { val = append(val, key); return true })
	return val
}

func (t *Set) Each(fn any) {
	assert.If(fn == nil, "[fn] should not be nil")

	vfn := reflect.ValueOf(fn)
	t.m.Range(func(key, value any) bool { _ = vfn.Call(funk.ListOf(reflect.ValueOf(key))); return true })
}
