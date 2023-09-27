package errors

import (
	"fmt"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func New(msg string) error {
	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Err:    parseToProto(msg),
			Caller: stack.Caller(1).String(),
		},
	}
}

func Debug(err error) {
	if generic.IsNil(err) {
		return
	}

	if _err, ok := err.(fmt.Stringer); ok {
		fmt.Println(_err.String())
		return
	}

	pretty.Println(err)
}

func UnwrapEach(err error, call func(e error) bool) {
	if err == nil {
		return
	}

	for {
		if !call(err) {
			return
		}

		err1, ok := err.(*ErrWrap)
		if !ok {
			return
		}

		err = err1.Unwrap()
	}
}

func WrapStack(err error) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(1).String(),
			Stack:  getStack(),
		},
	}
}

func WrapCaller(err error, skip ...int) error {
	if generic.IsNil(err) {
		return nil
	}

	var depth = 1
	if len(skip) > 0 {
		depth += skip[0]
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(depth).String(),
		},
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(1).String(),
			Tags:   map[string]string{"msg": fmt.Sprintf(format, args...)},
		},
	}
}

func Wrap(err error, msg string) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(1).String(),
			Tags:   map[string]string{"msg": msg},
		},
	}
}

func WrapTags(err error, tags map[string]string) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(1).String(),
			Tags:   tags,
		},
	}
}

func WrapKV(err error, key string, value any) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Caller: stack.Caller(1).String(),
			Tags:   map[string]string{key: fmt.Sprintf("%v", value)},
		},
	}
}
