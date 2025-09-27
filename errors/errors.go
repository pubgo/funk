package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pubgo/funk/v2/errors/errinter"
	"github.com/pubgo/funk/v2/proto/errorpb"
	"github.com/pubgo/funk/v2/stack"
)

func IfErr(err error, fn func(err error) error) error {
	if err == nil {
		return nil
	}

	return fn(err)
}

func New(msg string, tags ...Tag) error {
	return WrapCaller(&Err{Msg: msg, id: xid.New().String(), Tags: tags}, 1)
}

// NewFmt
// Deprecated: use Errorf instead
func NewFmt(msg string, args ...interface{}) error {
	return WrapCaller(&Err{Msg: fmt.Sprintf(msg, args...), id: xid.New().String()}, 1)
}

// Format
// Deprecated: use Errorf instead
func Format(msg string, args ...interface{}) error {
	return WrapCaller(&Err{Msg: fmt.Sprintf(msg, args...), id: xid.New().String()}, 1)
}

func Errorf(msg string, args ...interface{}) error {
	return WrapCaller(&Err{Msg: fmt.Sprintf(msg, args...), id: xid.New().String()}, 1)
}

func Parse(val interface{}) error { return errinter.ParseError(val) }
func Debug(err error)             { errinter.Debug(err) }

func Is(err, target error) bool {
	return errors.Is(err, target)
}
func Join(errs ...error) error { return errors.Join(errs...) }

func UnwrapEach(err error, call func(e error) bool) {
	if err == nil {
		return
	}

	for {
		if !call(err) {
			return
		}

		err1, ok := err.(ErrUnwrap)
		if !ok {
			return
		}

		err = err1.Unwrap()
	}
}

func AsA[T any](err error) (*T, bool) {
	var target T
	return &target, As(err, &target)
}

func As(err error, target any) bool {
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
	if err == nil {
		return nil
	}

	stack.PrintStack()
	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Stacks: lo.Map(getStack(), func(item *stack.Frame, index int) string { return item.String() }),
			Error:  MustProtoToAny(ParseErrToPb(err)),
		},
	}
}

func WrapCaller(err error, skip ...int) error {
	if err == nil {
		return nil
	}

	depth := 1
	if len(skip) > 0 {
		depth += skip[0]
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(depth).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
		},
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   Tags{T("msg", fmt.Sprintf(format, args...))}.ToMap(),
		},
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   Tags{T("msg", msg)}.ToMap(),
		},
	}
}

func WrapMapTag(err error, tags Maps) error {
	if err == nil {
		return nil
	}

	if tags == nil {
		return err
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   tags.Tags().ToMap(),
		},
	}
}

func WrapTag(err error, tags ...Tag) error {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   Tags(tags).ToMap(),
		},
	}
}

func WrapFn(err error, fn func() Tags) error {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   fn().ToMap(),
		},
	}
}

func WrapKV(err error, key string, value any, kvs ...any) error {
	if err == nil {
		return nil
	}

	var tags = Tags{T(key, value)}
	for i := 0; i < len(kvs); i += 2 {
		tags = append(tags, Tag{K: kvs[i].(string), V: kvs[i+1]})
	}

	return &ErrWrap{
		err: handleGrpcError(err),
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(ParseErrToPb(err)),
			Tags:   Tags{T(key, value)}.ToMap(),
		},
	}
}

func T(k string, v any) Tag {
	return Tag{K: k, V: v}
}

func MustTagsToAny(tags ...*errorpb.Tag) []*anypb.Any { return errinter.MustTagsToAny(tags...) }
func MustStructToAny(p map[string]any) *anypb.Any     { return errinter.MustStructToAny(p) }
func MustProtoToAny(p proto.Message) *anypb.Any       { return errinter.MustProtoToAny(p) }
func ParseErrToPb(err error) proto.Message            { return errinter.ParseErrToPb(err) }
func GetErrorId(err error) string                     { return errinter.GetErrorId(err) }
