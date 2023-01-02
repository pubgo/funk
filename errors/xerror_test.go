package errors

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	var err = fmt.Errorf("hello error")
	err = Wrap(err, "hello", "world")
	fmt.Printf("%q\n", err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
}
