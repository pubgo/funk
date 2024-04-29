package async

import (
	"context"
	"time"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
)

// Async 异步执行函数并同步等待
func Async[T any](fn func() (T, error)) *Future[T] {
	assert.If(fn == nil, "[Async] [fn] is nil")

	f := newFuture[T]()
	go func() {
		defer recovery.Recovery(func(err error) {
			err = errors.WrapKV(err, "fn_stack", stack.CallerWithFunc(fn).String())
			f.setErr(err)

			// 记录错误日志，以便跟踪
			logErr(fn, err)
		})

		t, e := fn()
		if e != nil {
			// 直接设置错误，上面的defer已经处理了错误包装
			f.setErr(e)
		} else {
			f.setOK(t)
		}
	}()
	return f
}

// GoSafe 安全并发处理
func GoSafe(fn func() error, cb ...func(err error)) {
	assert.If(fn == nil, "[GoSafe] [fn] is nil")

	go func() {
		err := try.Try(fn)
		if generic.IsNil(err) {
			return
		}

		if len(cb) > 0 && cb[0] != nil {
			logErr(cb[0], try.Try(func() error {
				cb[0](err)
				return nil
			}))
			return
		}

		logErr(fn, err)
	}()
}

// GoCtx 可取消并发处理
func GoCtx(fn func(ctx context.Context) error, parentCtx ...context.Context) context.CancelFunc {
	assert.If(fn == nil, "[GoCtx] [fn] is nil")

	pCtx := context.Background()
	if len(parentCtx) > 0 && parentCtx[0] != nil {
		pCtx = parentCtx[0]
	}

	ctx, cancel := context.WithCancel(pCtx)
	go func() {
		logErr(fn, try.Try(func() error { return fn(ctx) }))
	}()
	return cancel
}

// GoDelay 异步延迟处理, 默认10ms
func GoDelay(fn func() error, durations ...time.Duration) {
	assert.If(fn == nil, "[GoDelay] [fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	assert.If(dur == 0, "[dur] should not be 0")
	go func() { logErr(fn, try.Try(fn)) }()

	time.Sleep(dur)
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func() error) error {
	if fn == nil {
		return errors.New("[Timeout] [fn] is nil")
	}

	if dur <= 0 {
		return errors.New("[Timeout] [dur] should not be less than zero")
	}

	done := make(chan error)
	go func() {
		defer close(done)
		done <- try.Try(fn)
	}()

	timer := time.NewTimer(dur)
	defer timer.Stop()
	select {
	case <-timer.C:
		return context.DeadlineExceeded
	case ret := <-done:
		return ret
	}
}

func logErr(fn interface{}, err error) {
	if generic.IsNil(err) {
		return
	}

	logs.Err(err).
		Str("func", stack.CallerWithFunc(fn).String()).
		Str("err_stack", pretty.Sprint(err)).
		Msg(err.Error())
}
