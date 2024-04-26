package errors_test

import (
	"fmt"
	"testing"

	"github.com/rs/xid"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/version"
)

func TestCodeErr(t *testing.T) {
	err := errors.NewCodeErr(&errorpb.ErrCode{
		StatusCode: errorpb.Code_Aborted,
		Code:       100000,
		Name:       "hello.test.123",
		Message:    fmt.Sprintf("test error"),
	})

	err = errors.WrapMapTag(err, errors.Maps{
		"event":   "test event",
		"test123": 123,
	})

	err = errors.Wrap(err, "next error")
	err = errors.WrapTag(err, errors.T("test", "hello"))
	err = errors.Wrapf(err, "next error name=%s", "wrapf")

	err = errors.WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = errors.IfErr(err, func(err error) error {
		return errors.WrapMsg(err, &errorpb.ErrMsg{
			Msg: "this is if err msg",
		})
	})

	err = errors.WrapFn(err, func() errors.Tags {
		return errors.Tags{
			{"key", "map value"},
		}
	})

	err = errors.WrapTag(err, errors.T("name", "value"), errors.T("name1", "value"))
	err = errors.WrapTrace(err, &errorpb.ErrTrace{
		Version: version.Version(),
		Service: version.Project(),
		Id:      xid.New().String(),
	})

	err = errors.WrapStack(err)
	errors.Debug(err)

	var fff *errors.ErrCode
	t.Log(errors.As(err, &fff))
	t.Log(fff.Proto())
}
