package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
)

func ErrOf(fn func(err *Err)) *Err {
	var err = &Err{Caller: stack.Caller(1)}
	fn(err)
	return err
}

var _ Error = (*Err)(nil)

type Err struct {
	Caller *stack.Frame `json:"caller"`
	Err    error        `json:"err"`
	Msg    string       `json:"msg"`
	Detail string       `json:"detail"`
	Tags   Tags         `json:"tags"`
}

func (e Err) Unwrap() error {
	return e.Err
}

func (e Err) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any)
	if e.Msg != "" {
		data["msg"] = e.Msg
	}

	if e.Detail != "" {
		data["detail"] = e.Detail
	}

	if e.Caller != nil {
		data["caller"] = e.Caller
	}

	if e.Tags != nil {
		data["tags"] = e.Tags
	}

	if _err, ok := e.Err.(json.Marshaler); ok && _err != nil {
		data["cause"] = _err
	} else if e.Err != nil {
		data["err_msg"] = e.Err.Error()
		data["err_detail"] = fmt.Sprintf("%#v", e.Err)
	}
	return jjson.Marshal(data)
}

func (e Err) Error() string {
	if e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e Err) String() string {
	var buf = bytes.NewBuffer(nil)
	if e.Caller != nil {
		buf.WriteString(fmt.Sprintf("caller: %s\n", e.Caller.String()))
	}

	buf.WriteString(fmt.Sprintf("msg=%q detail=%q", e.Msg, e.Detail))
	if e.Tags != nil {
		buf.WriteString(fmt.Sprintf(" tags=%q", e.Tags))
	}

	if e.Err != nil {
		if _err, ok := e.Err.(fmt.Stringer); !ok {
			buf.WriteString(fmt.Sprintf(" err_msg=%q", e.Err.Error()))
			buf.WriteString(fmt.Sprintf(" err_detail=%s", pretty.Sprint(e.Err)))
		} else {
			buf.WriteString("\n====================================================================\n")
			buf.WriteString(_err.String())
		}
	}
	return buf.String()
}
