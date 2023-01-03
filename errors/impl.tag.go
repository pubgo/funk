package errors

import (
	"encoding/json"
	"fmt"
	"reflect"

	jjson "github.com/goccy/go-json"
)

var _ TagWrapper = (*errTagImpl)(nil)

type errTagImpl struct {
	err  error
	tags map[string]any
}

func (e errTagImpl) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = json.Marshal(e)
		if err != nil {
			_, _ = fmt.Fprint(s, err)
		} else {
			_, _ = fmt.Fprint(s, string(data))
		}
	default:
		_, _ = fmt.Fprint(s, e.String())
	}
}

func (e errTagImpl) Unwrap() error {
	return e.err
}

func (e errTagImpl) As(target interface{}) bool {
	if e == nil || target == nil {
		return false
	}

	var v = reflect.ValueOf(target)
	t1 := reflect.Indirect(v).Interface()
	if err, ok := t1.(TagWrapper); ok {
		v.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

func (e errTagImpl) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any)
	if e.tags != nil && len(e.tags) > 0 {
		data["tags"] = e.tags
	}

	if e.err != nil && !isNil(e.err) {
		data["err_msg"] = e.err.Error()
		data["err_detail"] = fmt.Sprintf("%#v", e.err)
	}
	return jjson.Marshal(data)
}

func (e errTagImpl) Error() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	return e.err.Error()
}

func (e errTagImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	return fmt.Sprintf("err=%q tags=%q", e.err.Error(), e.tags)
}

func (e errTagImpl) Tags() map[string]any {
	return e.tags
}
