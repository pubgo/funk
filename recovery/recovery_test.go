package recovery_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
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

	err := handler()
	if generic.IsNil(err) {
		t.Log(err)
	}
}

func TestResult(t *testing.T) {
	type A struct {
		A string
		B string
	}

	handler := func() (r result.Result[A]) {
		defer recovery.Err(&r.E)

		r = r.WithVal(A{A: "hello"})
		panic("ok")
	}

	t.Log(handler())
	t.Log(handler().OrElse(A{A: "error"}).A)
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
