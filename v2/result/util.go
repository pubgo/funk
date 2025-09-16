package result

import (
	"context"
	"fmt"
	"sync"

	"log/slog"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/log/logfields"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/v2/result/resultchecker"
)

var errFnIsNil = errors.New("[fn] is nil")

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errFnIsNil)
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			errors.Debug(gErr)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}

func try1[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errFnIsNil)
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			errors.Debug(gErr)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return
}

func errNilOrPanic(err error, events ...func(e *zerolog.Event)) {
	if err == nil {
		return
	}

	logErr(nil, err, events...)
	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}

func catchErr(r Error, setter ErrSetter, rawSetter *error, contexts ...context.Context) bool {
	if setter == nil && rawSetter == nil {
		errNilOrPanic(errors.Errorf("error setter is nil"))
	}

	if r.IsOK() {
		return false
	}

	var isErr = func() bool {
		if setter != nil {
			return setter.IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	var getErr = func() error {
		if setter != nil {
			return setter.GetErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil
	}

	var setErr = func(err error) {
		if setter != nil {
			setError(setter, err)
		}

		if rawSetter != nil {
			setError(ErrProxyOf(rawSetter), err)
		}
	}

	var ctx = context.Background()
	for i := range contexts {
		if contexts[i] == nil {
			continue
		}
		ctx = contexts[i]
		break
	}

	// err No checking, repeat setting
	if isErr() {
		err := getErr()
		log.Err(err, ctx).Msgf("error setter has already set the error, err=%s", err.Error())
	}

	var checkers = append(resultchecker.GetErrChecks(), resultchecker.GetCheckersFromCtx(ctx)...)
	var err = r.getErr()
	for _, fn := range checkers {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	setErr(errors.WrapCaller(err, 2))

	return true
}

func errRecovery(getErr func() error, callbacks ...func(err error) error) error {
	err := errors.Parse(recover())
	if err == nil {
		err = getErr()
	}

	if err == nil {
		return nil
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return nil
		}
	}

	debug.PrintStack()
	return err
}

func unwrapErr[T any](r Result[T], setter1 *error, setter2 ErrSetter, contexts ...context.Context) (T, error) {
	if setter1 == nil && setter2 == nil {
		errNilOrPanic(fmt.Errorf("error setter is nil"))
	}

	var ret = r.getValue()
	if r.IsOK() {
		return ret, nil
	}

	var ctx = context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}

	getErr := func() error {
		err := lo.FromPtr(setter1)
		if err == nil {
			err = setter2.GetErr()
		}
		return err
	}
	if preErr := getErr(); preErr != nil {
		log.Err(preErr, ctx).Msgf("error setter has already set the error, err=%v", preErr)
	}

	var err = r.getErr()
	var checkers = append(resultchecker.GetErrChecks(), resultchecker.GetCheckersFromCtx(ctx)...)
	for _, fn := range checkers {
		err = fn(ctx, err)
		if err == nil {
			return ret, nil
		}
	}

	return ret, err
}

func setError(setter ErrSetter, err error) {
	if err == nil {
		return
	}

	if setter == nil {
		errNilOrPanic(errors.Errorf("error setter is nil"))
		return
	}

	switch errSet := setter.(type) {
	case *Error:
		errSet.err = err
	case *ErrProxy:
		*errSet.err = err
	default:
		rv := reflect.ValueOf(setter)
		t := rv.Type()

		if !strings.Contains(t.String(), "Result[") {
			slog.Error("error setter type error, type is not Result",
				slog.String("type", fmt.Sprintf("%T", setter)),
				slog.String("type-string", t.String()),
				slog.String("stack", string(debug.Stack())),
			)
			return
		}

		ret := (*Result[any])(rv.UnsafePointer())
		ret.err = err
	}
}

func logErr(ctx context.Context, err error, events ...func(e *zerolog.Event)) {
	if err == nil {
		return
	}

	log.Error(ctx).
		Func(func(e *zerolog.Event) {
			e.Str(logfields.Module, "resultv2")
			e.Str(logfields.ErrorStack, string(debug.Stack()))
			e.Str(logfields.ErrorDetail, pretty().Sprint(err))
			e.Str(logfields.ErrorID, errors.GetErrorId(err))

			for _, fn := range events {
				fn(e)
			}
		}).
		Str(zerolog.ErrorFieldName, err.Error()).
		CallerSkipFrame(2).
		Msgf("%s\n%s\n", err.Error(), prototext.Format(errors.ParseErrToPb(err)))
}

var pretty = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(5)
	return printer
})
