package xtry

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
)

func TestTryErr(t *testing.T) {
	TryErr(func() result.Error {
		panic("ok")
	}).Do(func(err result.Error) {
		t.Log(err)
	})
}

func testFunc() (err error) {
	defer recovery.Err(&err, func(err errors.XErr) errors.XErr {
		return err.WrapF("test func")
	})
	assert.Must(errors.Err{Msg: "test error"})
	return
}

func TestTryVal(t *testing.T) {
	var v = TryVal(func() result.Result[*errors.Err] {
		return result.OK(&errors.Err{Msg: "ok"})
	})
	fmt.Println(v)
}
