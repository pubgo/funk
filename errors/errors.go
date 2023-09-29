package errors

import (
	"fmt"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"reflect"
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

func Is(err error, target any) bool {
	if target == nil {
		return err == target
	}

	var _, isTargetErr = target.(error)

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}

		if x, ok := err.(ErrEqual); ok && x.IsEqual(target) {
			return true
		}

		if x, ok := err.(ErrIs); ok && isTargetErr && x.Is(target.(error)) {
			return true
		}

		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func As(err error, target any) bool {
	if err == nil {
		return false
	}

	if target == nil {
		panic("errors: target cannot be nil")
	}

	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	targetType := typ.Elem()
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}

		if x, ok := err.(ErrAs); ok && x.As(target) {
			return true
		}

		err = Unwrap(err)
	}

	return false
}

func Unwrap(err error) error {
	u, ok := err.(ErrUnwrap)
	if !ok {
		return nil
	}

	return u.Unwrap()
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

func ParseToWrap(err error) *errorpb.ErrWrap {
	return parseErrToWrap(err)
}
