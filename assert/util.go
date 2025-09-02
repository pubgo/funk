package assert

import (
	"fmt"
	"log/slog"
	"runtime/debug"

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
		return fmt.Sprintf("%+v", msgAndArgs[0])
	}

	return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
}

func must(err error, messageArgs ...any) {
	if err == nil {
		return
	}

	message := messageFromMsgAndArgs(messageArgs...)
	if message == "" {
		message = err.Error()
	} else {
		message = fmt.Sprintf("msg:%s err:%s", message, err.Error())
	}

	if EnablePrintStack {
		slog.Error(message)
		debug.PrintStack()
	}

	panic(message)
}

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = fmt.Errorf("[fn] is nil")
		return
	}

	defer func() {
		if gErr != nil {
			gErr = fmt.Errorf("stack:%s, err:%w", stack.CallerWithFunc(fn).String(), gErr)
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
