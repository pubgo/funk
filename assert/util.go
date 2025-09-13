package assert

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"

	"github.com/k0kubun/pp/v3"
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

	return pp.Sprint(msgAndArgs...)
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

	slog.Error(message, slog.String("stack", string(debug.Stack())))
	panic(err)
}

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = fmt.Errorf("[fn] is nil")
		return
	}

	defer func() {
		if gErr != nil {
			gErr = fmt.Errorf("stack:%s, err:%w", reflect.TypeOf(fn).String(), gErr)
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			gErr = fmt.Errorf("%v", err)
		}
	}()

	gErr = fn()
	return
}
