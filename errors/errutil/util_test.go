package errutil

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/errors"
)

func TestJson(t *testing.T) {
	var err error = &errors.Err{
		Msg:    "this is msg",
		Detail: "this is detail",
		Tags:   errors.Tags{"tag": "hello"},
	}

	err = errors.Wrap(err, "this is next msg")
	fmt.Printf("%s", err)

	fmt.Println(string(JsonPretty(err)))
}
