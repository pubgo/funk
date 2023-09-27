package errors

import (
	"fmt"
	"testing"

	"github.com/rs/xid"

	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/version"
)

func TestCodeErr(t *testing.T) {
	var err = NewCodeErr(&errorpb.ErrCode{
		Code:    errorpb.Code_Aborted,
		BizCode: 100000,
		Name:    "hello.test.123",
		Reason:  fmt.Sprintf("test error"),
	})

	err = Wrap(err, "next error")
	err = Wrapf(err, "next error name=%s", "wrapf")

	err = WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = WrapTrace(err, &errorpb.ErrTrace{
		Version: version.Version(),
		Service: version.Project(),
		Id:      xid.New().String(),
	})

	err = WrapStack(err)
	Debug(err)
}
