package errutil

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/stack"
)

func TestJson(t *testing.T) {
	var err = errors.SimpleErr(func(err *errors.Err) {
		err.Err = fmt.Errorf("this is Err")
		err.Msg = "this is msg"
		err.Detail = "this is detail"
		err.Tags = errors.Tags{"tag": "hello"}
	})

	err = &errors.Err{Err: err, Msg: "this is next msg", Detail: "this is next detail", Caller: stack.Caller(0)}
	fmt.Printf("%s", err)

	fmt.Println(string(JsonPretty(err)))
}
