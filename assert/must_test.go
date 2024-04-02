package assert

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/stretchr/testify/assert"
)

type errBase struct {
	msg string
}

func panicErr() (*errBase, error) {
	return nil, fmt.Errorf("error")
}

func panicNoErr() (*errBase, error) {
	return &errBase{msg: "ok"}, nil
}

func TestPanicErr(t *testing.T) {
	is := assert.New(t)
	is.Panics(func() {
		ret := Must1(panicErr())
		fmt.Println(ret == nil)
	})

	is.NotPanics(func() {
		ret := Must1(panicNoErr())
		fmt.Println(ret.msg)
	})
}

func TestRespTest(t *testing.T) {
	defer func() {
		errors.Debug(errors.Parse(recover()))
	}()
	testPanic1(t)
}

func TestRespNext(t *testing.T) {
	testPanic1(t)
}

func testPanic1(t *testing.T) {
	// xerrImpl.Must(xerrImpl.New("ok"))
	Must(init1Next())
}

func init1Next() (err error) {
	Must(fmt.Errorf("test next"))
	return nil
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			Must(nil)
			return
		}()
	}
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				recover()
			}()

			panic("hello")
		}()
	}
}
