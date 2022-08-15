package syncx

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/xerr"
)

func Promise[T any](fn func(resolve func(T), reject func(error))) chan typex.Result[T] {
	assert.If(fn == nil, "[fn] is nil")

	var ch = make(chan typex.Result[T])
	go func() {
		defer close(ch)
		defer recovery.Recovery(func(err xerr.XErr) {
			ch <- typex.Err[T](err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
		})
		fn(func(t T) { ch <- typex.OK(t) }, func(err error) { ch <- typex.Err[T](err) })
	}()
	return ch
}
