package errutil

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestJson(t *testing.T) {
	var err = errors.SimpleErr(func(err *errors.Err) {
		err.Msg = "this is msg"
		err.Detail = "this is detail"
		err.Tags = errors.Tags{"tag": "hello"}
	})

	err = errors.Wrap(err, "this is next msg")
	fmt.Printf("%s", err)

	fmt.Println(string(JsonPretty(err)))
}
