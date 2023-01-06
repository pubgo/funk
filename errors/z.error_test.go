package errors

import (
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestFormat(t *testing.T) {
	var err = Err(New("hello error"))
	err = Wrap(err, "next error")
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = WrapTags(err, Map{"test": "hello"})
	err = WrapCode(err, codes.Canceled)
	err = WrapBizCode(err, "user.not_found")
	err = WrapStack(err)
	Debug(err)
	var ddd, _ = json.MarshalIndent(ParseResp(err), " ", "  ")
	fmt.Println(string(ddd))
}
