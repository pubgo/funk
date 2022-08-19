package strutil

import (
	"encoding"
	"fmt"
	"strings"
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
		b.WriteString(Stringify(v))
	}

	return b.String()
}

func Quote(s string) string {
	if strings.ContainsAny(s, " ") {
		return fmt.Sprintf("%q", s)
	}
	return s
}

func Stringify(v any) string {
	var s string
	switch vv := v.(type) {
	case string:
		s = vv
	case fmt.Stringer:
		s = vv.String()
	case error:
		s = vv.Error()
	case nil:
		return "null"
	case []byte:
		return string(vv)
	case encoding.TextMarshaler:
		vb, err := vv.MarshalText()
		if err != nil {
			return err.Error()
		}
		return string(vb)
	default:
		s = fmt.Sprint(v)
	}

	return Quote(s)
}
