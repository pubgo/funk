package vars

import (
	"expvar"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
)

func Float(name string) *expvar.Float {
	v := expvar.Get(name)
	if v == nil {
		return expvar.NewFloat(name)
	}
	return v.(*expvar.Float)
}

func Int(name string) *expvar.Int {
	v := expvar.Get(name)
	if v == nil {
		return expvar.NewInt(name)
	}
	return v.(*expvar.Int)
}

func String(name string) *expvar.String {
	v := expvar.Get(name)
	if v == nil {
		return expvar.NewString(name)
	}
	return v.(*expvar.String)
}

func Map(name string) *expvar.Map {
	v := expvar.Get(name)
	if v == nil {
		return expvar.NewMap(name)
	}
	return v.(*expvar.Map)
}

var _ json.Marshaler = (*Value)(nil)

type Value func() interface{}

func (f Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

func (f Value) Value() interface{} { return f() }

func (f Value) String() (r string) {
	var errStr = func(err any) string {
		ret := result.Wrap(json.Marshal(err))
		if ret.IsErr() {
			return pretty.Sprint(ret.Err())
		} else {
			return convert.B2S(ret.Unwrap())
		}
	}

	defer recovery.Recovery(func(err error) { r = errStr(err) })

	dt := f()
	switch dt.(type) {
	case nil:
		return "null"
	case string:
		return dt.(string)
	case []byte:
		return string(dt.([]byte))
	case fmt.Stringer:
		return dt.(fmt.Stringer).String()
	default:
		return errStr(dt)
	}
}

func Register(name string, value Value) {
	defer recovery.Exit()
	assert.If(Has(name), "name:%s already exists", name)
	expvar.Publish(name, value)
}

func RegisterValue(name string, data interface{}) {
	defer recovery.Exit()
	assert.If(Has(name), "name:%s already exists", name)
	expvar.Publish(name, Value(func() interface{} { return data }))
}

func Has(name string) bool {
	return expvar.Get(name) != nil
}

func Each(fn func(key string, val expvar.Var)) {
	expvar.Do(func(kv expvar.KeyValue) { fn(kv.Key, kv.Value) })
}
