package xerr

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	var err = fmt.Errorf("hello error")
	err = WrapXErr(err).WrapMeta("hello", "world")
	fmt.Printf("%q\n", err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
}
