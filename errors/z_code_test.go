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

	err = WrapMapTag(err, Maps{
		"event":   "test event",
		"test123": 123,
	})

	err = Wrap(err, "next error")
	err = WrapTag(err, T("test", "hello"))
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, &Err{Msg: "Err errors error", Tags: Tags{T("tags", "hello")}})

	err = WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = IfErr(err, func(err error) error {
		return WrapMsg(err, &errorpb.ErrMsg{
			Msg: "this is if err msg",
		})
	})

	err = WrapFn(err, func() Tags {
		return Tags{
			{"key", "map value"},
		}
	})

	err = WrapTag(err, T("name", "value"), T("name1", "value"))
	err = WrapTrace(err, &errorpb.ErrTrace{
		Version: version.Version(),
		Service: version.Project(),
		Id:      xid.New().String(),
	})

	err = WrapStack(err)
	Debug(err)

	var fff *ErrCode
	t.Log(As(err, &fff))
	t.Log(fff.Proto())
}
