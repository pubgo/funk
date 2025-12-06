package recovery_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func testExit1() {
	testExit()
}

func testExit() {
	defer recovery.Exit()

	assert.Must(fmt.Errorf("test"))
}

func TestExit(t *testing.T) {
	testExit1()
}

func TestErr(t *testing.T) {
	handler := func() (gErr error) {
		defer recovery.Err(&gErr)

		panic("ok")
	}
	t.Log("error:", handler())
}

func TestResult(t *testing.T) {
	type A struct {
		A string
		B string
	}

	handler := func() (r result.Result[A]) {
		defer result.Recovery(&r)

		r = r.WithValue(A{A: "hello"})
		panic("ok")
	}

	t.Log(handler())
	t.Log(handler().UnwrapOr(A{A: "error"}).A)
}

func TestName(t *testing.T) {
	defer recovery.DebugPrint()

	log.Print("test panic")
	hello()
}

func hello() {
	panic("hello")
}

func TestTesting(t *testing.T) {
	defer recovery.Testing(t)

	log.Print("test panic")
	hello()
}
