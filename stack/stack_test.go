package stack

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/pretty"
)

func init1() {
	init2()
}

func init2() {
	init3()
}

func init3() {
	pretty.Println(GetGORoot())
	pretty.Println(Callers(4))
	fmt.Println(Caller(0))
	fmt.Println(Caller(1))
	fmt.Println(Caller(2))
	fmt.Println(Caller(3))
	fmt.Println(Caller(20))
}

func TestCallerWithDepth(t *testing.T) {
	init1()
	fmt.Print("\n\n\n")
	init2()
	fmt.Print("\n\n\n")
	init3()
}
