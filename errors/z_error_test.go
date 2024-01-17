package errors_test

import (
	"fmt"
	"github.com/pubgo/funk/errors"
	"testing"

	"github.com/rs/xid"

	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/proto/testcodepb"
	"github.com/pubgo/funk/version"
)

func TestFormat(t *testing.T) {
	var err = errors.WrapCaller(fmt.Errorf("test error, err=%w", errors.New("hello error")))
	err = errors.Wrap(err, "next error")
	err = errors.WrapTag(err, errors.T("event", "test event"), errors.T("test123", 123), errors.T("test", "hello"))
	err = errors.Wrapf(err, "next error name=%s", "wrapf")

	err = errors.WrapCode(err, testcodepb.ErrCodeNotFound)
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

func TestNew(t *testing.T) {
	fmt.Printf("%s", errors.New("test error"))
}
