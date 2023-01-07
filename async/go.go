package async

import (
	"context"
	"time"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
)

// Async 通过chan的方式同步执行异步任务
func Async[T any](fn func() result.Result[T]) *Future[T] {
	assert.If(fn == nil, "[Async] [fn] is nil")

	var ch = newFuture[T]()
	go func() {
		ch.success(try.TryVal(fn))
	}()
	return ch
}

// GoSafe 安全并发处理
func GoSafe(fn func() result.Error, cb ...func(err result.Error)) {
	assert.If(fn == nil, "[GoSafe] [fn] is nil")

	go func() {
		var err = try.TryErr(fn)
		if err.IsNil() {
			return
		}

		if len(cb) > 0 && cb[0] != nil {
			logErr(cb[0], try.Try(func() { cb[0](err) }))
			return
		}

		logErr(fn, err)
	}()
}

// GoCtx 可取消并发处理
func GoCtx(fn func(ctx context.Context) result.Error) context.CancelFunc {
	assert.If(fn == nil, "[GoCtx] [fn] is nil")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		logErr(fn, try.TryErr(func() result.Error { return fn(ctx) }))
	}()
	return cancel
}

// GoDelay 异步延迟处理
func GoDelay(fn func() result.Error, durations ...time.Duration) {
	assert.If(fn == nil, "[GoDelay] [fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	assert.If(dur == 0, "[dur] should not be 0")
	go func() {
		logErr(fn, try.TryErr(fn))
	}()

	time.Sleep(dur)

	return
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func() result.Error) (gErr result.Error) {
	defer recovery.Recovery(func(err errors.XErr) {
		gErr = result.WithErr(err)
	})

	assert.If(fn == nil, "[Timeout] [fn] is nil")
	assert.If(dur <= 0, "[Timeout] [dur] should not be less than zero")

	var done = make(chan struct{})
	go func() {
		defer close(done)
		gErr = try.TryErr(fn)
	}()

	select {
	case <-time.After(dur):
		return result.WithErr(context.DeadlineExceeded)
	case <-done:
		return
	}
}

func logErr(fn interface{}, err result.Error) {
	if err.IsNil() {
		return
	}

	logs.Error(err.Unwrap(), err.Unwrap().Error(), "func", stack.CallerWithFunc(fn))
}
