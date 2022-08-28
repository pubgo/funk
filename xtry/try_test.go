package xtry

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
)

func TestTryErr(t *testing.T) {
	TryErr(func() result.Error {
		panic("ok")
	}).Do(func(err result.Error) {
		t.Log(err)
	})
}

func testFunc() (err error) {
	defer recovery.Err(&err, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("test func")
	})
	assert.Must(xerr.Err{Msg: "test error"})
	return
}

func TestTryVal(t *testing.T) {
	var v = TryVal(func() result.Result[*xerr.Err] {
		return result.OK(&xerr.Err{Msg: "ok"})
	})
	fmt.Println(v)
}
