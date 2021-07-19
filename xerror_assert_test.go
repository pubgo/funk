package xerror

import (
	"testing"
)

func TestAssertEqual(t *testing.T) {
	defer RespExit()

	AssertEqual("hello", 1)
}

func TestCheckNil(t *testing.T) {
	defer RespExit()

	var a *int
	Assert(a == nil, "ok")
}

func TestCheck(t *testing.T) {
	defer RespDebug("")

	Assert(true, "aaaa")
	Assert(false, "aaaa")
}
