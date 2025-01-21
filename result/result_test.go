package result_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
)

type hello struct {
	Name string `json:"name"`
}

func TestName(t *testing.T) {
	ok := result.OK(&hello{Name: "abc"})
	okBytes := result.Of(json.Marshal(&ok))
	data := string(okBytes.Expect("failed to encode json data"))
	t.Log(data)
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}

	var ok1 hello
	if err := json.Unmarshal([]byte(data), &ok1); err != nil {
		t.Fatal(err)
	}
	t.Log("ok", ok1.Name)
}

func TestResultDo(t *testing.T) {
	ok := result.OK(&hello{Name: "abc"})
	ok.Do(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	})
	ok.Do(func(v *hello) {
		assert.If(v.Name != "abc", "not match")
	})
}

func TestErrOf(t *testing.T) {
	errors.Debug(fn1().OnValue(func(tt string) error {
		t.Log(tt)
		return nil
	}))
}

func fn1() (r result.Result[string]) {
	fn3().ErrTo(&r)
	if r.IsErr() {
		return
	}

	var vv = fn2().ErrTo(&r)
	if r.IsErr() {
		return
	}

	return r.WithVal(vv)
}

func fn2() (r result.Result[string]) {
	fn3().ErrTo(&r, func(err error) error {
		return errors.Wrap(err, "test error")
	})

	if r.IsErr() {
		return
	}

	return r.WithVal("ok")
}

func fn3() result.Error {
	return result.ErrOf(fmt.Errorf("error"))
}
