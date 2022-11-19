package strutil

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var jsonMarshaller = &protojson.MarshalOptions{EmitUnpopulated: true}

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
	case nil:
		s = "null"
	case string:
		s = vv
	case fmt.Stringer:
		s = vv.String()
	case error:
		s = vv.Error()
	case []byte:
		s = string(vv)
	case json.Marshaler:
		vb, err := vv.MarshalJSON()
		if err != nil {
			s = err.Error()
		} else {
			s = string(vb)
		}
	case encoding.TextMarshaler:
		vb, err := vv.MarshalText()
		if err != nil {
			s = err.Error()
		} else {
			s = string(vb)
		}
	case proto.Message:
		var dt, err = jsonMarshaller.Marshal(vv)
		if err != nil {
			s = err.Error()
		} else {
			s = string(dt)
		}
	default:
		var dt, err = json.Marshal(v)
		if err != nil {
			s = err.Error()
		} else {
			s = string(dt)
		}
	}
	return Quote(s)
}
