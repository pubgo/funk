package errors

import (
	"encoding/json"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestFormat(t *testing.T) {
	var err = New("hello error")
	err = Wrap(err, "next error")
	err = WrapTags(err, Map{"test": "hello"})
	err = WrapCode(err, codes.Canceled)
	err = WrapBizCode(err, "user.not_found")
	err = WrapStack(err)
	Debug(err)
	var ddd, _ = json.MarshalIndent(ParseResp(err), "  ", "  ")
	t.Log(string(ddd))
}
