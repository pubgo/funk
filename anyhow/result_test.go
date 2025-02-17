package anyhow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/anyhow"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/recovery"
)

type hello struct {
	Name string `json:"name"`
}

func TestName(t *testing.T) {
	defer recovery.DebugPrint()
	ok := &hello{Name: "abc"}
	okBytes := anyhow.Wrap(json.Marshal(&ok))
	data := string(okBytes.Expect("failed to encode json data"))
	t.Log(data)
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}
}

func TestResultDo(t *testing.T) {
	ok := anyhow.OK(&hello{Name: "abc"})
	ok.OnValue(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	}).OnValue(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	}).OnErr(func(err error) {})
}

func TestErrOf(t *testing.T) {
	anyhow.RegisterErrCheck(func(err error) error {
		return errors.Wrap(err, "global err check")
	})

	errors.Debug(fn1().OnValue(func(tt string) {
		t.Log(tt)
	}).GetErr())
}

func fn1() (r anyhow.Result[string]) {
	var ctx = log.UpdateEventCtx(context.Background(), log.Map{"test": "ok"})
	fn3().ErrTo(&r.Err, log.RecordErr(ctx))
	if r.IsErr() {
		return
	}

	var vv = fn2().ErrTo(&r.Err)
	if r.IsErr() {
		return
	}

	return r.WithVal(vv)
}

func fn2() (r anyhow.Result[string]) {
	fn3().ErrTo(&r.Err, func(err error) error {
		return errors.Wrap(err, "test error")
	})

	if r.IsErr() {
		return
	}

	return r.WithVal("ok")
}

func fn3() anyhow.Error {
	return anyhow.ErrOf(fmt.Errorf("error"))
}
