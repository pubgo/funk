package result_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
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

	assert.Assert(err1().Err() == nil, "failed to check CatchTo")
	assert.Assert(err1().Err().Error() != "test error", "error not match")
}

func err1() (gErr result.Result[any]) {
	ret := result.Err[any](fmt.Errorf("test error"))
	if ret.CatchTo(&gErr.E) {
		return
	}
	return
}
