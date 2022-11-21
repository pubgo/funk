package stack

import (
	"fmt"
	"testing"
)

func init1() {
	init2()
}

func init2() {
	init3()
}

func init3() {
	fmt.Println(Callers(4))
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
