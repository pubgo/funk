package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/pubgo/funk/logger"
	_ "github.com/rs/xid"
)

type jsonLogger struct {
	io.Writer
}

func (l *jsonLogger) Info(level uint, msg string, tags ...logger.Tagger) {
	//TODO implement me
	panic("implement me")
}

func (l *jsonLogger) Error(err error, msg string, tags ...logger.Tagger) {
	//TODO implement me
	panic("implement me")
}

// NewJSONLogger returns a Logger that encodes keyvals to the Writer as a
// single JSON object. Each log event produces no more than one call to
// w.Write. The passed Writer must be safe for concurrent use by multiple
// goroutines if the returned Logger will be used concurrently.
func NewJSONLogger(w io.Writer) logger.Logger {
	return &jsonLogger{w}
}

func (l *jsonLogger) Log(tags ...logger.Tagger) error {
	m := make(map[string]interface{}, len(tags))
	for i := range tags {
		merge(m, tags[i].Key(), tags[i].Value())
	}
	enc := json.NewEncoder(l.Writer)
	enc.SetEscapeHTML(false)
	return enc.Encode(m)
}

func merge(dst map[string]interface{}, k, v interface{}) {
	var key string
	switch x := k.(type) {
	case string:
		key = x
	case fmt.Stringer:
		key = safeString(x)
	default:
		key = fmt.Sprint(x)
	}

	// We want json.Marshaler and encoding.TextMarshaller to take priority over
	// err.Error() and v.String(). But json.Marshall (called later) does that by
	// default so we force a no-op if it's one of those 2 case.
	switch x := v.(type) {
	case json.Marshaler:
	case encoding.TextMarshaler:
	case error:
		v = safeError(x)
	case fmt.Stringer:
		v = safeString(x)
	}

	dst[key] = v
}

func safeString(str fmt.Stringer) (s string) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(str); v.Kind() == reflect.Ptr && v.IsNil() {
				s = "NULL"
			} else {
				s = fmt.Sprintf("PANIC in String method: %v", panicVal)
			}
		}
	}()
	s = str.String()
	return
}

func safeError(err error) (s interface{}) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				s = nil
			} else {
				s = fmt.Sprintf("PANIC in Error method: %v", panicVal)
			}
		}
	}()
	s = err.Error()
	return
}
