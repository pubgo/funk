package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/stack"
)

func New(format string, a ...interface{}) error {
	return &errImpl{
		err:    fmt.Errorf(format, a...),
		caller: stack.Caller(1),
	}
}

var _ IMsgWrap = (*errMsgImpl)(nil)

type errMsgImpl struct {
	*errImpl
	msg string
}

func (e errMsgImpl) Msg() string {
	return e.msg
}

func (e errMsgImpl) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["msg"] = e.msg
	return jjson.Marshal(data)
}

func (e errMsgImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("msg"), e.msg))
	buf.WriteString(e.errImpl.String())
	return buf.String()
}
