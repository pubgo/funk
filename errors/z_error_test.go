package errors

import (
	"fmt"
	"testing"

	"github.com/rs/xid"

	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/proto/testcodepb"
	"github.com/pubgo/funk/version"
)

func TestFormat(t *testing.T) {
	var err = WrapCaller(fmt.Errorf("test error, err=%w", New("hello error")))
	err = Wrap(err, "next error")
	err = WrapTags(err, Tags{"event": "test event", "test123": 123, "test": "hello"})
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, &Err{Msg: "Err errors error", Tags: map[string]any{"tags": "hello"}})

	err = WrapCode(err, testcodepb.ErrCodeNotFound)
	err = WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = IfErr(err, func(err error) error {
		return WrapMsg(err, &errorpb.ErrMsg{
			Msg: "this is if err msg",
		})
	})

	err = WrapFn(err, func(tag Tags) {
		tag.Set("key", "map value")
	})

	err = WrapTag(err, T("name", "value"), T("name", "value"))
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
