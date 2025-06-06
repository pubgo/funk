package errors_test

import (
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/proto/testcodepb"
)

func TestErrCode(t *testing.T) {
	err1 := errors.Wrap(errors.NewCodeErr(testcodepb.TestErrCodeDbConn), "hello")
	rr := errors.Is(err1, errors.NewCodeErr(testcodepb.TestErrCodeDbConn))
	if !rr {
		t.Fatal("not match")
	}

	t.Log(errors.As(err1, testcodepb.TestErrCodeDbConn))
}
