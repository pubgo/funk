package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"google.golang.org/grpc/codes"
)

var _ ICodeWrap = (*errCodeImpl)(nil)

type errCodeImpl struct {
	*errImpl
	code codes.Code
}

func (e errCodeImpl) Code() codes.Code {
	return e.code
}

func (e errCodeImpl) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["code"] = e.code.String()
	return jjson.Marshal(data)
}

func (e errCodeImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("code"), e.code.String()))
	buf.WriteString(e.errImpl.String())
	return buf.String()
}
