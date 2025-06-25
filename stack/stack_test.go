package stack_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
	"github.com/stretchr/testify/assert"
)

func init1() {
	init2()
}

func init2() {
	init3()
}

func init3() {
	pretty.Println("GetGORoot", stack.GetGORoot())
	pretty.Println("Callers", stack.Callers(4))
	fmt.Println(stack.Caller(0))
	fmt.Println(stack.Caller(1))
	fmt.Println(stack.Caller(2))
	fmt.Println(stack.Caller(3))
	fmt.Println(stack.Caller(20))
}

func TestCallerWithDepth(t *testing.T) {
	fmt.Println(stack.Caller(0).String())
	init1()
	fmt.Print("\n\n\n")
	init2()
	fmt.Print("\n\n\n")
	init3()
}

func TestCallType(t *testing.T) {
	assert.Equal(t,
		"github.com/pubgo/funk/errors",
		stack.CallerWithType(reflect.TypeOf(errors.ErrMsg{})).Pkg,
	)
}
