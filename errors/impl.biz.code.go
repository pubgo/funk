package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
)

var _ IBizCodeWrap = (*errBizCodeImpl)(nil)

type errBizCodeImpl struct {
	*errImpl
	bizCode string
}

func (e errBizCodeImpl) BizCode() string {
	return e.bizCode
}

func (e errBizCodeImpl) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["biz_code"] = e.bizCode
	return jjson.Marshal(data)
}

func (e errBizCodeImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("biz"), e.bizCode))
	buf.WriteString(e.errImpl.String())
	return buf.String()
}
