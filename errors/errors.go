package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/stack"
)

func IfErr(err error, fn func(err error) error) error {
	if err == nil {
		return nil
	}

	return fn(err)
}

func New(msg string, tags ...Tags) error {
	return WrapCaller(newSimpleErr(&Err{Msg: msg, id: NewErrorId(), Tags: lo.FirstOrEmpty(tags)}), 1)
}

func Errorf(msg string, args ...any) error {
	return WrapCaller(&Err{Msg: fmt.Sprintf(msg, args...), id: NewErrorId()}, 1)
}

func Is(err, target error) bool { return errors.Is(err, target) }
func Join(errs ...error) error  { return errors.Join(errs...) }
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
	u, ok := err.(ErrUnwrapper)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func WrapStack(err error) error {
	if err == nil {
		return nil
	}

	stack.Print()
	return newErrWrapStack(err, Tags{"msg": err.Error()})
}

func WrapTagsCaller(err error, tags Tags, skip ...int) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, tags, lo.FirstOrEmpty(skip))
}

func WrapCaller(err error, skip ...int) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, Tags{"msg": err.Error()}, lo.FirstOrEmpty(skip))
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, Tags{"msg": fmt.Sprintf(format, args...)})
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, Tags{"msg": msg})
}

func WrapTags(err error, tags ...Tags) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, mergeTags(tags...))
}

func WrapFn(err error, tagsFn ...func() Tags) error {
	if err == nil {
		return nil
	}

	tags := lo.Map(tagsFn, func(item func() Tags, index int) Tags { return item() })
	return newErrWrap(err, mergeTags(tags...))
}

func WrapKV(err error, key string, value any) error {
	if err == nil {
		return nil
	}

	return newErrWrap(err, Tags{key: value})
}

func JsonPrint(err error) []byte {
	if err == nil {
		return nil
	}

	data, err := json.Marshal(err)
	if err != nil {
		slog.Error("failed to marshal error", "err", err)
		panic(fmt.Errorf("failed to marshal error, err=%w", err))
	}
	return data
}

func DebugPrint(err error) {
	if err == nil {
		return
	}

	if _err, ok := err.(fmt.Stringer); ok {
		_, _ = fmt.Fprintln(os.Stderr, _err.String())
		return
	}

	debugPretty().Println(err)
}

func GetErrorId(err error) string { return getErrorId(err) }
