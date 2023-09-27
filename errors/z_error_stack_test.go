package errors

import (
	"testing"
)

func TestStack(t *testing.T) {
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	err = WrapStack(err)
	err = Wrapf(err, "next error name=%s", "wrapf")
	Debug(err)
}
