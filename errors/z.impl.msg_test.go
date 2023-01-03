package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	err := Err{
		Err:    fmt.Errorf("hello"),
		Msg:    "test",
		Detail: "test test",
	}

	var dd, _ = json.Marshal(err)
	t.Log(string(dd))
	t.Log(err.String())
	t.Log(err.Error())
}
