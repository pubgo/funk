package error1

import (
	"fmt"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

type Tags map[string]any

type Errors interface {
	Errors() []error
	Append(err error) error
}

type ErrUnwrap interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

func WrapCode(err error, code *errorpb.ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if generic.IsNil(code) {
		return nil
	}

	var wrap = ParseErr(err)
	wrap.wrap.Err = &errorpb.ErrWrap_Code{Code: code}
	wrap.wrap.Caller = stack.Caller(1).String()
	return wrap
}

func WrapTrace(err error, trace *errorpb.ErrTrace) error {
	if generic.IsNil(err) {
		return nil
	}

	if generic.IsNil(trace) {
		return nil
	}

	var wrap = ParseErr(err)
	wrap.wrap.Err = &errorpb.ErrWrap_Trace{Trace: trace}
	wrap.wrap.Caller = stack.Caller(1).String()
	return wrap
}

func Wrap(err error) error {
	if generic.IsNil(err) {
		return nil
	}

	var wrap = ParseErr(err)
	wrap.wrap.Caller = stack.Caller(1).String()
	return wrap
}

func WrapTags(err error, tags ...*errorpb.Tag) error {
	if generic.IsNil(err) {
		return nil
	}

	var wrap = ParseErr(err)
	wrap.wrap.Tags = tags
	wrap.wrap.Caller = stack.Caller(1).String()
	return wrap
}

func WrapMsg(err error, msg string, args ...any) error {
	if generic.IsNil(err) {
		return nil
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	var wrap = ParseErr(err)
	wrap.wrap.Tags = append(wrap.wrap.Tags, &errorpb.Tag{Key: "msg", Value: msg})
	wrap.wrap.Caller = stack.Caller(1).String()
	return wrap
}
