package logfmt

import (
	"encoding/json"
	"strconv"
	"time"

	"strings"

	_ "github.com/gookit/goutil/stdutil"
	_ "github.com/gookit/goutil/strutil"
	_ "github.com/gookit/slog"
)

func Flatten(kvs ...interface{}) string {
	if len(kvs) == 0 {
		return ""
	}

	var b strings.Builder
	for i := 0; i < len(kvs); i += 2 {
		if i > 0 {
			b.WriteRune(' ')
		}

		k := kvs[i]
		var v interface{}
		if i+1 < len(kvs) {
			v = kvs[i+1]
		} else {
			v = ""
		}
		b.WriteString(Stringify(k))
		b.WriteRune('=')
		b.WriteString(Quote(Stringify(v)))
	}

	return b.String()
}

func Quote(s string) string {
	if strings.ContainsAny(s, " ") {
		return strconv.Quote(s)
	}
	return s
}

func Stringify(v any) string {
	if v == nil {
		return "null"
	}

	switch vv := v.(type) {
	case string:
		return vv
	case int:
		return strconv.Itoa(vv)
	case int8:
		return strconv.Itoa(int(vv))
	case int16:
		return strconv.Itoa(int(vv))
	case int32: // same as `rune`
		return strconv.Itoa(int(vv))
	case int64:
		return strconv.FormatInt(vv, 10)
	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case uint64:
		return strconv.FormatUint(vv, 10)
	case float32:
		return strconv.FormatFloat(float64(vv), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vv, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(vv)
	case []byte:
		return string(vv)
	case time.Duration:
		return vv.String()
	case json.Number:
		return vv.String()
	case complex64:
		return `"` + strconv.FormatComplex(complex128(vv), 'f', -1, 64) + `"`
	case complex128:
		return `"` + strconv.FormatComplex(vv, 'f', -1, 128) + `"`
	default:
		return ""
	}
}
