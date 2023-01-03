package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

func New(format string, a ...interface{}) Error {
	return &errImpl{
		err:    fmt.Errorf(format, a...),
		tags:   make(map[string]any),
		msg:    fmt.Sprintf(format, a...),
		status: codes.Unknown,
		caller: stack.Caller(1),
	}
}

type errImpl struct {
	err    error
	msg    string
	caller *stack.Frame
}

func (t *errImpl) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any)
	if t.msg != "" {
		data["msg"] = t.msg
	}

	data["status"] = t.status

	if t.bizCode != "" {
		data["biz_code"] = t.bizCode
	}

	if t.caller != nil {
		data["caller"] = t.caller
	}

	data["stack_trace"] = t.stackTrace

	if t.tags != nil && len(t.tags) > 0 {
		data["tags"] = t.tags
	}

	if t.err != nil {
		data["err_msg"] = t.err.Error()
		data["err_detail"] = fmt.Sprintf("%#v", t.err)
	}
	return jjson.Marshal(data)
}

func (t *errImpl) Code() codes.Code {
	return t.status
}

func (t *errImpl) String() string { return t.debugString() }

func (t *errImpl) Unwrap() error { return t.err }

func (t *errImpl) _p(buf *strings.Builder, xrr *errImpl) {
	if xrr == nil || isNil(xrr) {
		return
	}

	buf.WriteString("========================================================================================================================\n")
	if xrr.err != nil {
		buf.WriteString(fmt.Sprintf(" %s]: %q\n", color.Red.P("error"), xrr.err.Error()))
	}

	if xrr.msg != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("msg"), xrr.msg))
	}

	if xrr.status != 0 {
		buf.WriteString(fmt.Sprintf("  %s]: %q\n", color.Green.P("code"), xrr.status.String()))
	}

	if xrr.bizCode != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("biz"), xrr.bizCode))
	}

	if xrr.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("caller"), xrr.caller.String()))
	}

	for i := range xrr.stackTrace {
		if xrr.stackTrace[i].IsRuntime() {
			continue
		}

		buf.WriteString(fmt.Sprintf(" %s]: %s\n", color.Yellow.P("stack"), xrr.stackTrace[i].String()))
	}

	if len(xrr.tags) > 0 {
		var dd, err = json.MarshalIndent(xrr.tags, "  ", "  ")
		if err != nil {
			panic(err)
		}
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("tags"), string(dd)))
	}

	t._p(buf, trans(xrr.err))
}

func (t *errImpl) debugString() string {
	if t == nil || t.err == nil {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString("\n")
	t._p(buf, t)
	buf.WriteString("========================================================================================================================\n\n")
	return buf.String()
}

func (t *errImpl) Is(err error) bool {
	if t == nil || t.err == nil || err == nil {
		return false
	}

	switch _err := err.(type) {
	case *errImpl:
		return _err == t || _err.err == t.err
	case error:
		return t.err == _err
	default:
		return false
	}
}

func (t *errImpl) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = json.Marshal(t)
		if err != nil {
			_, _ = fmt.Fprint(s, err)
		} else {
			_, _ = fmt.Fprint(s, string(data))
		}
	default:
		_, _ = fmt.Fprint(s, t.debugString())
	}
}

func (t *errImpl) As(target interface{}) bool {
	if t == nil || target == nil {
		return false
	}

	var v = reflect.ValueOf(target)
	t1 := reflect.Indirect(v).Interface()
	if err, ok := t1.(*errImpl); ok {
		v.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

// Error
func (t *errImpl) Error() string {
	if t == nil || isNil(t.err) {
		return ""
	}

	return t.err.Error()
}
