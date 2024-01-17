package result_test

import (
	"encoding/json"
	"github.com/pubgo/funk/result"
	"testing"
)

type hello struct {
	Name string `json:"name"`
	Jj   func()
}

func TestName(t *testing.T) {
	var ok = result.OK(&hello{Name: "abc"})
	okBytes := result.Of(json.Marshal(&ok))
	data := string(okBytes.Expect("failed to encode json data"))
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}

	var ok1 result.Result[*hello]
	if err := json.Unmarshal([]byte(data), &ok1); err != nil {
		t.Fatal(err)
	}
	t.Log("ok", ok1.Unwrap().Name)
}
