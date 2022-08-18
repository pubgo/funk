package xerr

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var err = fmt.Errorf("hello error")
	err = WrapXErr(err)
	fmt.Printf("%q", err)
}
