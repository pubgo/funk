package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
)

var _ Error = (*ErrCode)(nil)
var _ fmt.Formatter = (*ErrCode)(nil)

type ErrCode struct {
	err error
	pb  *errorpb.ErrCode
}

func (t *ErrCode) Kind() string {
	return "code"
}

func (t *ErrCode) Format(f fmt.State, verb rune) {
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

func (t *ErrCode) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))

	if t.pb.Code != 0 {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCode, t.code.String()))
	}

	if t.pb.reason != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorReason, t.reason))
	}

	if t.pb.status != "" {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStatus, t.status))
	}

	if len(t.pb.tags) > 0 {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorTags, t.tags))
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, t.caller.String()))
	}

	errStringify(buf, t.err)

	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["tags"] = t.tags
	data["status"] = t.status
	data["code"] = t.code.String()
	data["reason"] = t.reason
	return jjson.Marshal(data)
}

func (t *ErrCode) getData() map[string]any {
	var data = make(map[string]any)
	data["kind"] = t.Kind()
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

func (t *ErrCode) Unwrap() error { return t.err }

// Error
func (t *ErrCode) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
