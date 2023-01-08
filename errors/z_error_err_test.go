package errors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pubgo/funk/stack"
)

func TestErr(t *testing.T) {
	var err = ErrOf(func(err *Err) {
		err.Err = fmt.Errorf("this is Err")
		err.Msg = "this is msg"
		err.Detail = "this is detail"
		err.Tags = Tags{"tag": "hello"}
	})

	err = &Err{Err: err, Msg: "this is nest msg", Detail: "this is nest detail", Caller: stack.Caller(0)}
	fmt.Printf("%s\n", err)

	var data, _ = json.MarshalIndent(err, "  ", "  ")
	fmt.Println(string(data))
}
