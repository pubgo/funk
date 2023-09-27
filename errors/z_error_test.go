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
	err = Wrapf(err, "next error name=%s", "wrapf")

	err = WrapCode(err, testcodepb.ErrCodeNotFound)
	err = WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = WrapTrace(err, &errorpb.ErrTrace{
		Version: version.Version(),
		Service: version.Project(),
		Id:      xid.New().String(),
	})

	err = WrapStack(err)
	fmt.Println(err.(*ErrWrap).String())
	Debug(err)
}
