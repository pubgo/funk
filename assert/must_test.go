package assert_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	assert1 "github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/debugs"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
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
		errors.DebugPrint(errparser.Parse(recover()))
	}()
	assert1.Must(init1Next())
}

func TestRespNext(t *testing.T) {
	is := assert.New(t)
	is.Panics(func() {
		_ = init1Next()
	})
}

func init1Next() (err error) {
	assert1.Must(fmt.Errorf("test next"))
	return nil
}

func TestDebugMode(t *testing.T) {
	is := assert.New(t)
	wasDebug := debugs.IsDebug()
	t.Cleanup(func() {
		if wasDebug {
			debugs.SetEnabled()
		} else {
			debugs.SetDisabled()
		}
	})

	assert1.Exit(debugs.Enabled.Set("true"))
	is.Panics(func() {
		assert1.Must(fmt.Errorf("test next"))
	})
}

func BenchmarkNoPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() (err error) {
			assert1.Must(nil)
			return err
		}()
	}
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				_ = recover()
			}()

			panic("hello")
		}()
	}
}
