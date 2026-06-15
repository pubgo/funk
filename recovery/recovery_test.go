package recovery_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

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
	var code int
	recovery.SetExitFn(func(c int) { code = c })
	t.Cleanup(func() { recovery.SetExitFn(nil) })

	testExit1()
	require.Equal(t, 1, code)
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

func TestName(_ *testing.T) {
	defer recovery.DebugPrint()

	log.Print("test panic")
	hello()
}

func hello() {
	panic("hello")
}

func TestTesting(t *testing.T) {
	var fatalErr error
	recovery.SetTestingFatalFn(func(_ *testing.T, err error) { fatalErr = err })
	t.Cleanup(func() { recovery.SetTestingFatalFn(nil) })

	func() {
		defer recovery.Testing(t)

		log.Print("test panic")
		hello()
	}()

	require.EqualError(t, fatalErr, "hello")
}
