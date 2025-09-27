package vars

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"go.uber.org/atomic"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
)

var mux sync.Mutex

func Bool(name string) *atomic.Bool {
	return Any(name, atomic.NewBool(false))
}

func Float(name string) *atomic.Float64 {
	return Any(name, atomic.NewFloat64(0))
}

func Int(name string) *atomic.Int64 {
	return Any(name, atomic.NewInt64(0))
}

func String(name string) *atomic.String {
	return Any(name, atomic.NewString(""))
}

func Duration(name string) *atomic.Duration {
	return Any(name, atomic.NewDuration(0))
}

func Time(name string) *atomic.Time {
	return Any(name, atomic.NewTime(time.Now()))
}

func Error(name string) *atomic.Error {
	return Any(name, atomic.NewError(nil))
}

var _ json.Marshaler = (*Func)(nil)

type Func func() any

func (f Func) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

func (f Func) Value() any { return f() }

func (f Func) String() (r string) {
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
	case json.Marshaler:
		return jsonStr(dt)
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

func Register(name string, value Func) {
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
