package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	var err = Err{Err: fmt.Errorf("this is Err"), Msg: "this is msg", Detail: "this is detail"}
	err = Err{Err: err, Msg: "this is nest msg", Detail: "this is nest detail"}
	fmt.Println(err.String())

	var data, _ = json.MarshalIndent(err, "  ", "  ")
	fmt.Println(string(data))
}
