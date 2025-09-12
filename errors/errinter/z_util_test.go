package errinter_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/proto/errorpb"
)

func TestGetErrorId(t *testing.T) {
	err := errors.New("test error id")
	pbTxt := prototext.Format(errinter.ParseErrToPb(err))
	errId := errinter.GetErrorId(err)
	assert.NotEmpty(t, errId)
	assert.Contains(t, pbTxt, errId)
}

func TestParseError(t *testing.T) {
	assert.Equal(t, errinter.ParseError(nil), nil)
	assert.Equal(t, errinter.ParseError("hello").Error(), "hello")
	assert.Equal(t, errinter.ParseError(fmt.Errorf("err test")).Error(), "err test")

	errMsg := errinter.ParseError(prototext.Format(errinter.MustProtoToAny(&errorpb.ErrMsg{Msg: "test msg error"}))).Error()
	assert.Equal(t,
		strings.ReplaceAll(errMsg, "  ", " "),
		"[type.googleapis.com/errors.ErrMsg]: {\n msg: \"test msg error\"\n}\n",
	)

	ctx := context.WithValue(context.Background(), "hello", "world")
	assert.Equal(t,
		errinter.ParseError(ctx).Error(),
		"&context.valueCtx{\n  Context: context.backgroundCtx{},\n  key:     \"hello\",\n  val:     \"world\",\n}",
	)
}
