package async

import (
	"context"
	"time"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
)

// Async 通过chan的方式同步执行异步任务
func Async[T any](fn func() (T, error)) *result.Future[T] {
	assert.If(fn == nil, "[Async] [fn] is nil")

	var ch = result.NewFuture[T]()
	go func() {
		ch.Err(try.Try(func() error {
			var v, err = fn()
			if err != nil {
				return err
			}

			ch.OK(v)
			return nil
		}))
	}()
	return ch
}

// GoSafe 安全并发处理
func GoSafe(fn func() error, cb ...func(err error)) {
	assert.If(fn == nil, "[GoSafe] [fn] is nil")

	go func() {
		var err error
		try.WithErr(&err, fn)
		if errors.IsNil(err) {
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
func GoCtx(fn func(ctx context.Context) error) context.CancelFunc {
	assert.If(fn == nil, "[GoCtx] [fn] is nil")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		logErr(fn, try.Try(func() error { return fn(ctx) }))
	}()
	return cancel
}

// GoDelay 异步延迟处理
func GoDelay(fn func() error, durations ...time.Duration) {
	assert.If(fn == nil, "[GoDelay] [fn] is nil")

	dur := time.Millisecond * 10
	if len(durations) > 0 {
		dur = durations[0]
	}

	assert.If(dur == 0, "[dur] should not be 0")
	go func() { logErr(fn, try.Try(fn)) }()

	time.Sleep(dur)

	return
}

// Timeout 超时处理
func Timeout(dur time.Duration, fn func() error) (gErr result.Error) {
	assert.If(fn == nil, "[Timeout] [fn] is nil")
	assert.If(dur <= 0, "[Timeout] [dur] should not be less than zero")

	var done = make(chan struct{})
	go func() {
		defer close(done)
		gErr = gErr.WithErr(try.Try(fn))
	}()

	select {
	case <-time.After(dur):
		return result.WithErr(context.DeadlineExceeded)
	case <-done:
		return
	}
}

func logErr(fn interface{}, err error) {
	if errors.IsNil(err) {
		return
	}

	logs.Err(err).
		Str("func", stack.CallerWithFunc(fn).String()).
		Msg(err.Error())
}
