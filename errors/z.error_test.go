package errors

import (
	"testing"

	"google.golang.org/grpc/codes"
)

func TestFormat(t *testing.T) {
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	err = WrapFn(err, func(xrr XError) {
		xrr.AddMsg("new error msg")
		xrr.AddTag("test123", 123)
	})
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = WrapTags(err, Tags{"test": "hello"})
	err = WrapCode(err, codes.Canceled)
	err = WrapBizCode(err, "user.not_found")
	err = WrapStack(err)
	Debug(err)
}
