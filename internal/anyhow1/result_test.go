package anyhow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	anyhow "github.com/pubgo/funk/internal/anyhow1"
	"github.com/pubgo/funk/internal/anyhow1/aherrcheck"
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
	}).OnErr(func(err error) error {
		return err
	})
}

func TestErrOf(t *testing.T) {
	var ctx = log.UpdateEventCtx(context.Background(), log.Map{"test": "ok"})
	aherrcheck.RegisterErrCheck(log.RecordErr())

	var err anyhow.Error
	if fn1().ErrTo(&err, ctx) {
		errors.Debug(err.GetErr())
	}
}

func fn1() (r anyhow.Result[string]) {
	if fn3().ErrTo(&r.Err) {
		return
	}

	var vv = fn2()
	if vv.ErrTo(&r.Err) {
		return
	}

	return vv
}

func fn2() (r anyhow.Result[string]) {
	if fn3().OnErr(func(err error) error {
		return errors.Wrap(err, "test error")
	}).ErrTo(&r.Err) {
		return
	}

	return r.WithVal("ok")
}

func fn3() anyhow.Error {
	return anyhow.ErrOf(fmt.Errorf("error"))
}
