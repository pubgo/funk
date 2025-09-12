package vars

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/xid"
	"go.uber.org/atomic"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
)

var mux sync.Mutex

func Bool(name string) *atomic.Bool {
	return Any(name, atomic.NewBool(false))
}

func Float(name string) *expvar.Float {
	mux.Lock()
	defer mux.Unlock()

	v := expvar.Get(name)
	if v == nil {
		return expvar.NewFloat(name)
	}
	return v.(*expvar.Float)
}

func Int(name string) *expvar.Int {
	mux.Lock()
	defer mux.Unlock()

	v := expvar.Get(name)
	if v == nil {
		return expvar.NewInt(name)
	}
	return v.(*expvar.Int)
}

func String(name string) *expvar.String {
	mux.Lock()
	defer mux.Unlock()

	v := expvar.Get(name)
	if v == nil {
		return expvar.NewString(name)
	}
	return v.(*expvar.String)
}

func Map(name string) *expvar.Map {
	mux.Lock()
	defer mux.Unlock()

	v := expvar.Get(name)
	if v == nil {
		return expvar.NewMap(name)
	}
	return v.(*expvar.Map)
}

var _ json.Marshaler = (*Value)(nil)

type Value func() any

func (f Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

func (f Value) Value() any { return f() }

func (f Value) String() (r string) {
	return toString(f())
}

func errToString(err error) string {
	if err == nil {
		return "null"
	}

	return strconv.Quote(fmt.Sprintf("err:%s detail:%#v", err.Error(), err))
}
func toString(dt any) (r string) {
	var jsonStr = func(data any) string {
		ret, err := json.Marshal(data)
		if err != nil {
			return errToString(err)
		} else {
			return string(ret)
		}
	}

	defer recovery.Recovery(func(err error) { r = jsonStr(err) })

	switch dt := dt.(type) {
	case nil:
		return "null"
	case string:
		return strconv.Quote(dt)
	case []byte:
		return strconv.Quote(string(dt))
	case fmt.Stringer:
		return strconv.Quote(dt.String())
	case error:
		return errToString(dt)
	default:
		return slog.AnyValue(dt).String()
	}
}

func Any[T any](name string, v T) T {
	mux.Lock()
	defer mux.Unlock()

	vv := expvar.Get(name)
	if vv != nil {
		return vv.(*anyValue).v.(T)
	}

	expvar.Publish(name, &anyValue{v: v})
	return v
}

var _ expvar.Var = (*anyValue)(nil)

type anyValue struct {
	v any
}

func (a anyValue) String() string {
	return toString(a.v)
}

func (a anyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.v)
}

func (a anyValue) Value() any { return a.v }

func Register(name string, value Value) {
	assert.If(Has(name), "name:%s already exists", name)
	expvar.Publish(name, value)
}

func Has(name string) bool {
	return expvar.Get(name) != nil
}

func Each(fn func(key string, val expvar.Var)) {
	expvar.Do(func(kv expvar.KeyValue) { fn(kv.Key, kv.Value) })
}

func UniqueName(names ...string) string {
	return strings.Join(append(names, xid.New().String()), "_")
}
