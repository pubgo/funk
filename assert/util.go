package assert

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/pubgo/funk/log/logfields"
	"github.com/samber/lo"
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

func logErr(err error, message string, attrs ...slog.Attr) {
	if err == nil {
		return
	}

	attrs = append(attrs,
		slog.String(logfields.Module, "assert"),
		slog.String(logfields.Error, err.Error()),
		slog.String(logfields.ErrorStack, string(debug.Stack())),
		slog.String(logfields.ErrorDetail, pretty().Sprint(err)),
	)
	slog.Error(message, lo.ToAnySlice(attrs)...)
}

func must(err error, messageArgs ...any) {
	if err == nil {
		return
	}

	message := messageFromMsgAndArgs(messageArgs...)
	if message == "" {
		message = err.Error()
	} else {
		message = fmt.Sprintf("msg:%v err:%s", message, err.Error())
	}

	logErr(err, message, slog.Bool("panic", true))
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
