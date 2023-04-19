package error1

import (
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/proto/testcodepb"
	"testing"
)

var err1 = errors.Append(errors.New("raw error"), errors.New("raw error"), errors.New("raw error"))
var err2 = New("error").WithTag("module", "errors")

func TestName(t *testing.T) {
	err3 := Wrap(err1)
	err3 = WrapMsg(err3, "this is message")
	err3 = WrapCode(err3, testcodepb.ErrCodeNotFound)
	err3 = WrapTags(err3, &errorpb.Tag{Key: "new_tag", Value: "tag"})

	var err = ParseErr(err3)
	t.Log(err.Error())
	t.Log(err.String())
	ss, _ := err.MarshalJSON()
	t.Log(string(ss))
	log.Err(err).RawJSON("ddd", ss).Msg("error")
	t.Log(GetError(err).String())
}

func TestName1(t *testing.T) {
	err3 := Wrap(err2)
	err3 = WrapMsg(err3, "this is message")
	err3 = WrapCode(err3, testcodepb.ErrCodeNotFound)
	err3 = WrapTags(err3, &errorpb.Tag{Key: "new_tag", Value: "tag"})

	var err = ParseErr(err3)
	t.Log(err.Error())
	t.Log(err.String())
	t.Log(err.Proto().String())
	ss, _ := err.MarshalJSON()
	t.Log(string(ss))
	t.Log(GetError(err).String())
}
