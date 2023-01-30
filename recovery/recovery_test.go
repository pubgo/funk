package recovery

import (
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/result"
)

func TestErr(t *testing.T) {
	var handler = func() (gErr error) {
		defer Err(&gErr)

		panic("ok")
	}

	var err = handler()
	if errors.IsNil(err) {
		t.Log(err)
	}
}

func TestResult(t *testing.T) {
	type A struct {
		A string
		B string
	}

	var handler = func() (r result.Result[A]) {
		defer Result(&r)

		r = r.WithVal(A{A: "hello"})
		panic("ok")
	}

	t.Log(handler())
	t.Log(handler().OrElse(A{A: "error"}).A)
}

func TestName(t *testing.T) {
	defer DebugPrint()

	log.Print("test panic")
	hello()
}

func hello() {
	panic("hello")
}

func TestDump(t *testing.T) {
	defer Dump()
	hello()
}
