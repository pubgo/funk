package errutil_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errutil"
)

func TestJson(t *testing.T) {
	var err error = &errors.Err{
		Msg:    "this is msg",
		Detail: "this is detail",
		Tags:   errors.Tags{errors.T("tag", "hello")},
	}

	err = errors.Wrap(err, "this is next msg")
	fmt.Printf("%s", err)

	fmt.Println(string(errutil.JsonPretty(err)))
}
