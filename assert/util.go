package assert

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	"github.com/pubgo/funk/log/logfields"
	"github.com/pubgo/funk/stack"
)

func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}

	if len(msgAndArgs) == 1 {
		if msgAsStr, ok := msgAndArgs[0].(string); ok {
			return msgAsStr
		}
	}

	return pretty().Sprint(msgAndArgs...)
}

var assetFile = stack.Caller(0)

func logErr(err error, message string, attrs ...slog.Attr) {
	if err == nil {
		return
	}

	traces := lo.Filter(stack.Trace(), func(item *stack.Frame, index int) bool {
		return !item.IsRuntime() && item.Pkg != assetFile.Pkg
	})
	attrs = append(attrs,
		slog.String(logfields.Module, "assert"),
		slog.String(logfields.Error, err.Error()),
		slog.Any(logfields.ErrorStack, lo.Map(traces, func(item *stack.Frame, index int) string { return item.String() })),
		slog.String(logfields.ErrorDetail, fmt.Sprintf("%v", err)),
	)
	slog.Error(message, lo.ToAnySlice(attrs)...)
}

func must(err error, messageArgs ...any) {
	if err == nil {
		return
	}

	var attrs = []slog.Attr{slog.Bool("panic", true)}
	if v, ok := lo.ErrorsAs[interface {
		ID() string
		Error() string
	}](err); ok && v != nil {
		attrs = append(attrs, slog.String(logfields.ErrorID, v.ID()))
	}

	message := messageFromMsgAndArgs(messageArgs...)
	if message == "" {
		if v, ok := lo.ErrorsAs[interface {
			Proto() proto.Message
			Error() string
		}](err); ok && v != nil {
			message = fmt.Sprintf("%s\n%s", err.Error(), prototext.Format(v.Proto()))
		} else {
			message = err.Error()
		}
	} else {
		message = fmt.Sprintf("msg:%v err:%s", message, err.Error())
	}

	logErr(err, message, attrs...)
	panic(err)
}

var pretty = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(5)
	return printer
})

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = fmt.Errorf("assert: [fn] is nil")
		logErr(gErr, gErr.Error())
		debug.PrintStack()
		return
	}

	defer func() {
		if err := recover(); err != nil {
			gErr = fmt.Errorf("%v", err)
			logErr(gErr, gErr.Error())
			debug.PrintStack()
		}

		if gErr != nil {
			gErr = fmt.Errorf("stack:%s, err:%w", reflect.TypeOf(fn).String(), gErr)
		}
	}()

	gErr = fn()
	return
}
