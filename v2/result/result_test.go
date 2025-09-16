package result_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/result/resultchecker"
)

type hello struct {
	Name string `json:"name"`
}

func TestName(t *testing.T) {
	defer recovery.DebugPrint()
	ok := &hello{Name: "abc"}
	okBytes := result.Wrap(json.Marshal(&ok))
	data := string(okBytes.Expect("failed to encode json data"))
	t.Log(data)
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}
}

func TestResultDo(t *testing.T) {
	ok := result.OK(&hello{Name: "abc"})
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
	resultchecker.RegisterErrCheck(log.RecordErr())

	var err result.Error
	if fn1().Catch(&err, ctx) {
		errors.Debug(err.GetErr())
	}
}

func fn1() (r result.Result[string]) {
	if fn3().Catch(&r) {
		return
	}

	val := fn2().Unwrap(&r)
	if r.IsErr() {
		return
	}

	return r.WithValue(val)
}

func fn2() (r result.Result[string]) {
	fn3().
		InspectErr(func(err error) {
			log.Err(err).Msg("test error")
		}).
		Catch(&r)
	if r.IsErr() {
		return
	}

	return r.WithValue("ok")
}

func fn3() result.Error {
	return result.
		ErrOf(fmt.Errorf("error test, this is error")).
		InspectErr(func(err error) {
			log.Err(err).Msg("ddd")
		}).
		Log()
}
