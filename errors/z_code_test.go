package errors_test

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/stretchr/testify/assert"
)

func TestWrapCaller(t *testing.T) {
	err := fmt.Errorf("test")
	var ff = func() error {
		return errors.WrapCaller(err, 1)
	}

	assert.Contains(t, fmt.Sprint(ff()), "z_code_test.go:20 TestWrapCaller")
}

func TestCodeErr(t *testing.T) {
	err := errors.NewCodeErr(&errorpb.ErrCode{
		StatusCode: errorpb.Code_Aborted,
		Code:       100000,
		Name:       "hello.test.123",
		Message:    "test error",
	})

	err = errors.WrapMapTag(err, errors.Maps{
		"event":   "test event",
		"test123": 123,
	})

	err = errors.Wrap(err, "next error")
	err = errors.WrapTag(err, errors.T("test", "hello"))
	err = errors.Wrapf(err, "next error name=%s", "wrapf")

	err = errors.WrapMsg(err, &errorpb.ErrMsg{
		Msg: "this is msg",
	})

	err = errors.IfErr(err, func(err error) error {
		return errors.WrapMsg(err, &errorpb.ErrMsg{
			Msg: "this is if err msg",
		})
	})

	err = errors.WrapFn(err, func() errors.Tags {
		return errors.Tags{
			{"key", "map value"},
		}
	})

	err = errors.WrapTag(err, errors.T("name", "value"), errors.T("name1", "value"))

	err = errors.WrapStack(err)
	errors.Debug(err)

	var fff *errors.ErrCode
	t.Log(errors.As(err, &fff))
	t.Log(fff.Proto())
}
