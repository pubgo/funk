package errors

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestFormat(t *testing.T) {
	var err = New("hello error")
	err = Parse(err).WithTag("test", "hello").WithCode(codes.Canceled).WithBizCode("user.not_found")
	fmt.Printf("%q\n", err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
}
