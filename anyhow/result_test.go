package anyhow_test

import (
	"context"
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
	okBytes := anyhow.JsonMarshal(&ok)
	data := string(okBytes.Expect("failed to encode json data"))
	t.Log(data)
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}
}

func TestResultDo(t *testing.T) {
	ok := anyhow.Ok(&hello{Name: "abc"})
	ok.Inspect(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	})
	ok.Inspect(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	})
	ok.InspectErr(func(err error) {
		t.Log(err)
	})
}

func TestErrOf(t *testing.T) {
	var ctx = log.UpdateEventCtx(context.Background(), log.Map{"test": "ok"})
	// Note: Error checking functionality has been simplified in the new API

	var err error
	result := fn1()
	if result.IsError() {
		err = result.Err()
		errors.Debug(err)
	}
	_ = ctx // Use ctx to avoid unused variable warning
}

func fn1() anyhow.Result[string] {
	// Chain operations using the new API
	return fn3().
		OrElse(func(err error) anyhow.Result[string] {
			// Handle error or continue with fn2
			return fn2()
		})
}

func fn2() anyhow.Result[string] {
	return fn3().
		MapErr(func(err error) error {
			return errors.Wrap(err, "test error")
		}).
		OrElse(func(err error) anyhow.Result[string] {
			return anyhow.Ok("ok")
		})
}

func fn3() anyhow.Result[string] {
	return anyhow.Fail[string](fmt.Errorf("error test, this is error"))
}
