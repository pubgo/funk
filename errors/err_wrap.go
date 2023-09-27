package errors

import (
	"bytes"
	"fmt"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/proto/errorpb"
)

type ErrWrap struct {
	pb *errorpb.ErrWrap
}

func (e *ErrWrap) Unwrap() error {
	if e.pb.Wrap == nil {
		return nil
	}

	return &ErrWrap{pb: e.pb.Wrap}
}

func (e *ErrWrap) Error() string {
	if e.pb.Wrap == nil {
		return e.pb.String()
	}

	return (&ErrWrap{pb: e.pb.Wrap}).Error()
}

func (e *ErrWrap) Proto() *errorpb.ErrWrap { return e.pb }

func (e *ErrWrap) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString("===============================================================\n")
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, e.pb.Caller))
	for i := range e.pb.Tags {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorTags, e.pb.Tags[i]))
	}

	for i := range e.pb.Stack {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStack, e.pb.Stack[i]))
	}

	if e.pb.Err != nil {
		buf.WriteString(e.pb.Err.String())
	}

	if e.pb.Wrap != nil {
		buf.WriteString(e.pb.Wrap.String())
	}

	return buf.String()
}
