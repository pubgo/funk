package stack

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/kr/pretty"
)

func init1() {
	init2()
}

func init2() {
	init3()
}

func init3() {
	pretty.Log(Callers(4))
	pp.Println(Callers(4))
	fmt.Println(CallerWithDepth(0))
	fmt.Println(CallerWithDepth(1))
	fmt.Println(CallerWithDepth(2))
	fmt.Println(CallerWithDepth(3))
	fmt.Println(CallerWithDepth(20))
}

func TestCallerWithDepth(t *testing.T) {
	init1()
	fmt.Print("\n\n\n")
	init2()
	fmt.Print("\n\n\n")
	init3()
}
