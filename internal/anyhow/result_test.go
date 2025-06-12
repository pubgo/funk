package anyhow_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/internal/anyhow"
	"github.com/pubgo/funk/internal/anyhow/aherrcheck"
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
	ok.Inspect(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	}).Inspect(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	})
	ok.InspectErr(func(err error) {
		t.Log(err)
	})
}

func TestErrOf(t *testing.T) {
	var ctx = log.UpdateEventCtx(context.Background(), log.Map{"test": "ok"})
	aherrcheck.RegisterErrCheck(log.RecordErr())

	var err anyhow.Error
	if fn1().CatchErr(&err, ctx) {
		errors.Debug(err.GetErr())
	}
}

func fn1() (r anyhow.Result[string]) {
	if fn3().Catch(&r.Err) {
		return
	}

	var vv = fn2()
	if vv.Catch(&r.Err) {
		return
	}

	return vv
}

func fn2() (r anyhow.Result[string]) {
	ret := fn3().Map(func(err error) error {
		return errors.Wrap(err, "test error")
	})

	if ret.Catch(&r.Err) {
		return
	}

	return r.SetWithValue("ok")
}

func fn3() anyhow.Error {
	return anyhow.ErrOf(fmt.Errorf("error test, this is error")).
		Inspect(func(err error) {
			log.Err(err).Msg("ddd")
		}).
		InspectLog(func(evt *log.Event) {
			evt.Msg("test log")
		}).
		RecordLog()
}
