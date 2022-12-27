package stack

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/kr/pretty"
	"github.com/rs/zerolog/log"
)

func init1() {
	init2()
}

func init2() {
	init3()
}

func init3() {
	pp.Println(GetGORoot())
	pretty.Log(Callers(4))
	pp.Println(Callers(4))
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

func TestName(t *testing.T) {
	t.Log(pp.Sprint(CallerWithFunc(log.Info)))
	t.Log(CallerWithFunc(log.Info).Short())
}
