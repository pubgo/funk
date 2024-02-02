package errors_test

import (
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/proto/testcodepb"
	"testing"
)

func TestErrCode(t *testing.T) {
	var err1 = errors.Wrap(errors.NewCodeErr(testcodepb.ErrCodeDbConn), "hello")
	rr := errors.Is(err1, errors.NewCodeErr(testcodepb.ErrCodeDbConn))
	if !rr {
		t.Fatal("not match")
	}

	t.Log(errors.As(err1, testcodepb.ErrCodeDbConn))
}
