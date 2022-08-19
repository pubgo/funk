package xtry

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
)

func testFunc() (err error) {
	defer recovery.Err(&err, func(err xerr.XErr) xerr.XErr {
		return err.WrapF("test func")
	})
	assert.Must(xerr.Err{Msg: "test error"})
	return
}

func TestTryCatch(t *testing.T) {
	TryCatch(
		func() error {
			panic("ok")
		},
		func(err xerr.XErr) {
			err.DebugPrint()
		},
	)
}

func TestTryVal(t *testing.T) {
	var v = TryCatch1(func() result.Result[*xerr.Err] {
		return result.OK(&xerr.Err{Msg: "ok"})
	}, func(err xerr.XErr) {
		err.DebugPrint()
	})
	fmt.Println(v)
}
