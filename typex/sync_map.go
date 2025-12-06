package typex

import (
	"reflect"
	"sync"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
)

var NotFound = new(struct{})

type SyncMap struct {
	data sync.Map
}

func (t *SyncMap) Each(fn any) (err error) {
	defer recovery.Err(&err)

	assert.If(fn == nil, "[fn] should not be nil")

	vfn := reflect.ValueOf(fn)
	onlyKey := reflect.TypeOf(fn).NumIn() == 1
	t.data.Range(func(key, value any) bool {
		if onlyKey {
			_ = vfn.Call(funk.ListOf(reflect.ValueOf(key)))
			return true
		}

		_ = vfn.Call(funk.ListOf(reflect.ValueOf(key), reflect.ValueOf(value)))
		return true
	})

	return nil
}

func (t *SyncMap) Map(fn func(val any) any) {
	t.data.Range(func(key, value any) bool {
		t.data.Store(key, fn(value))
		return true
	})
}

func (t *SyncMap) MapTo(data any) (err error) {
	defer recovery.Err(&err)

	vd := reflect.ValueOf(data)
	if vd.Kind() == reflect.Ptr {
		vd = vd.Elem()
		vd.Set(reflect.MakeMap(vd.Type()))
	}

	// var data = make(map[string]int); MapTo(data)
	// var data map[string]int; MapTo(&data)
	assert.If(!vd.IsValid() || vd.IsNil(), "[data] type error")

	t.data.Range(func(key, value any) bool {
		vd.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		return true
	})

	return nil
}

func (t *SyncMap) Set(key, value any) {
	t.data.Store(key, value)
}

func (t *SyncMap) Get(key any) any {
	value, ok := t.data.Load(key)
	if ok {
		return value
	}

	return NotFound
}

func (t *SyncMap) LoadAndDelete(key any) (value any, ok bool) {
	return t.data.LoadAndDelete(key)
}
func (t *SyncMap) Load(key any) (value any, ok bool) { return t.data.Load(key) }
func (t *SyncMap) Range(f func(key, value any) bool) { t.data.Range(f) }
func (t *SyncMap) Delete(key any)                    { t.data.Delete(key) }
func (t *SyncMap) Has(key any) (ok bool)             { _, ok = t.data.Load(key); return ok }
