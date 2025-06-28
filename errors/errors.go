package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func IfErr(err error, fn func(err error) error) error {
	if err == nil {
		return nil
	}

	return fn(err)
}

func New(msg string) error {
	return WrapCaller(&Err{Msg: msg, id: xid.New().String()}, 1)
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

func NewTags(msg string, tags ...Tag) error {
	return WrapCaller(&Err{Msg: msg, Tags: tags, id: xid.New().String()}, 1)
}

func Parse(val interface{}) error {
	return parseError(val)
}

func Debug(err error) {
	if err == nil {
		return
	}

	err = parseError(err)
	if _err, ok := err.(fmt.Stringer); ok {
		fmt.Println(_err.String())
		return
	}

	pretty.SetDefaultMaxDepth(20)
	pretty.Println(err)
}

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

	debug.PrintStack()
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

func MustProtoToAny(p proto.Message) *anypb.Any {
	switch p := p.(type) {
	case nil:
		return nil
	case *anypb.Any:
		return p
	}

	pb, err := anypb.New(p)
	if err != nil {
		log.Err(err).Str("protobuf", prototext.Format(p)).Msgf("failed to encode protobuf message to any")
		return nil
	} else {
		return pb
	}
}

func ParseErrToPb(err error) proto.Message {
	switch err1 := err.(type) {
	case nil:
		return nil
	case ErrorProto:
		return err1.Proto()
	case GRPCStatus:
		return err1.GRPCStatus().Proto()
	case proto.Message:
		return err1
	default:
		return &errorpb.ErrMsg{Msg: err.Error(), Detail: fmt.Sprintf("%v", err)}
	}
}

func GetErrorId(err error) string {
	if err == nil {
		return ""
	}

	for err != nil {
		if v, ok := err.(Error); ok {
			return v.ID()
		}

		err = Unwrap(err)
	}

	return ""
}
