package assert_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	assert1 "github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
)

func init() {
	slog.SetDefault(slog.New(log.NewSlog(log.GetLogger(assert1.Name))))
}

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
		ret := assert1.Must1(panicErr())
		fmt.Println(ret == nil)
	})

	is.NotPanics(func() {
		ret := assert1.Must1(panicNoErr())
		fmt.Println(ret.msg)
	})
}

func TestRespTest(t *testing.T) {
	defer func() {
		errors.Debug(errors.Parse(recover()))
	}()
	assert1.Must(init1Next())
}

func TestRespNext(t *testing.T) {
	assert1.Must(init1Next())
}

func init1Next() (err error) {
	assert1.Must(fmt.Errorf("test next"))
	return nil
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			assert1.Must(nil)
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
