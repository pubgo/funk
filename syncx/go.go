package syncx

import (
	"context"
	"time"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
	"github.com/pubgo/funk/xtry"
)

// Async 通过chan的方式同步执行异步任务
func Async[T any](fn func() result.Result[T]) *Future[T] {
	assert.If(fn == nil, "[Async] [fn] is nil")

	var ch = newFuture[T]()
	go func() {
		ch.success(xtry.TryVal(fn))
	}()
	return ch
}

// GoSafe 安全并发处理
func GoSafe(fn func() result.Error, cb ...func(err error)) {
	assert.If(fn == nil, "[GoSafe] [fn] is nil")

	go func() {
		defer recovery.Recovery(func(err xerr.XErr) {
			if len(cb) > 0 && cb[0] != nil {
				logErr(cb[0], xtry.Try(func() { cb[0](err) }))
				return
			}

			logErr(fn, result.WithErr(err))
		})

		fn().Must()
	}()
}

// GoCtx 可取消并发处理
func GoCtx(fn func(ctx context.Context), cb ...func(err error)) context.CancelFunc {
	assert.If(fn == nil, "[GoCtx] [fn] is nil")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer recovery.Recovery(func(err xerr.XErr) {
			if len(cb) > 0 && cb[0] != nil {
				logErr(cb[0], xtry.Try(func() { cb[0](err) }))
				return
			}

			logErr(fn, result.WithErr(err))
		})

		fn(ctx)
	}()

	return cancel
}

// GoDelay 异步延迟处理
func GoDelay(fn func(), durations ...time.Duration) {
	assert.If(fn == nil, "[GoDelay] [fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	assert.If(dur == 0, "[dur] should not be 0")

	go func() {
		logErr(fn, xtry.Try(fn))
	}()

	time.Sleep(dur)

	return
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func()) (gErr error) {
	defer recovery.Err(&gErr)

	if dur <= 0 {
		panic("[Timeout] [dur] should not be less than zero")
	}

	assert.If(fn == nil, "[Timeout] [fn] is nil")

	var done = make(chan struct{})

	go func() {
		defer close(done)
		defer recovery.Err(&gErr)

		fn()
	}()

	select {
	case <-time.After(dur):
		return context.DeadlineExceeded
	case <-done:
		return
	}
}

func logErr(fn interface{}, err result.Error) {
	if err.IsNil() {
		return
	}

	logs.Error(err.Unwrap(), err.Unwrap().Error(), "func", utils.CallerWithFunc(fn))
}
