package errors

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	err = WrapTag(err,
		T("event", "test event"),
		T("test123", 123),
		T("test", "hello"),
	)

	err = WrapStack(err)
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, &Err{Msg: "Err errors error", Tags: Tags{T("tags", "hello")}})
	Debug(err)
}
