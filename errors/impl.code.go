package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

var _ ErrCode = (*errCodeImpl)(nil)
var _ fmt.Formatter = (*errCodeImpl)(nil)

type errCodeImpl struct {
	err    error
	caller *stack.Frame

	reason string
	status string
	code   errorpb.Code
	tags   map[string]string
}

func (t *errCodeImpl) SetStatus(status string) ErrCode {
	t.status = status
	return t
}

func (t *errCodeImpl) Kind() string {
	return "code"
}

func (t *errCodeImpl) Status() string {
	return t.status
}

func (t *errCodeImpl) SetErr(err error) ErrCode {
	t.err = err
	return t
}

func (t *errCodeImpl) Tags() map[string]string {
	return t.tags
}

func (t *errCodeImpl) AddTag(key string, val string) ErrCode {
	if t.tags == nil {
		t.tags = make(map[string]string)
	}
	t.tags[key] = val
	return t
}

func (t *errCodeImpl) Code() errorpb.Code {
	return t.code
}

func (t *errCodeImpl) Reason() string {
	return t.reason
}

func (t *errCodeImpl) SetCode(code errorpb.Code) ErrCode {
	t.code = code
	return t
}

func (t *errCodeImpl) SetReason(reason string) ErrCode {
	t.reason = reason
	return t
}

func (t *errCodeImpl) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = t.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, t.String())
	}
}

func (t *errCodeImpl) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))

	if t.code != 0 {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCode, t.code.String()))
	}

	if t.reason != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorReason, t.reason))
	}

	if t.status != "" {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStatus, t.status))
	}

	if len(t.tags) > 0 {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorTags, t.tags))
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, t.caller.String()))
	}

	errStringify(buf, t.err)

	return buf.String()
}

func (t *errCodeImpl) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["tags"] = t.tags
	data["status"] = t.status
	data["code"] = t.code.String()
	data["reason"] = t.reason
	return jjson.Marshal(data)
}

func (t *errCodeImpl) getData() map[string]any {
	var data = make(map[string]any)
	if t.caller != nil {
		data["caller"] = t.caller
	}

	var mm = errJsonify(t.err)
	if mm != nil {
		for k, v := range mm {
			data[k] = v
		}
	}

	return data
}

func (t *errCodeImpl) Unwrap() error { return t.err }

// Error
func (t *errCodeImpl) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
