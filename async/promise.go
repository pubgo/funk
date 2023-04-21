package async

import (
	"sync"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/stack"
)

func Promise[T any](fn func(resolve func(T), reject func(err error))) *Future[T] {
	assert.If(fn == nil, "[fn] is nil")

	var f = newFuture[T]()
	go func() {
		defer recovery.Recovery(func(err error) {
			err = errors.WrapKV(err, "fn", stack.CallerWithFunc(fn).String())
			f.setErr(err)
		})

		fn(func(t T) { f.setOK(t) }, func(err error) { f.setErr(err) })
	}()
	return f
}

func Group[T any](do func(async func(func() (T, error))) error) *Iterator[T] {
	assert.If(do == nil, "[Async] [fn] is nil")

	var rr = iteratorOf[T]()
	go func() {
		var wg sync.WaitGroup
		defer rr.setDone()
		defer wg.Wait()
		rr.setErr(do(func(f func() (T, error)) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer recovery.Recovery(func(err error) {
					err = errors.WrapKV(err, "fn_stack", stack.CallerWithFunc(do).String())
					rr.setErr(err)
				})

				var t, e = f()
				if e == nil {
					rr.setValue(t)
				} else {
					rr.setErr(e)
				}
			}()
		}))
	}()

	return rr
}

func Yield[T any](do func(yield func(T)) error) *Iterator[T] {
	var dd = iteratorOf[T]()
	go func() {
		defer dd.setDone()
		defer recovery.Recovery(func(err error) {
			err = errors.WrapTag(err, "fn_stack", stack.CallerWithFunc(do).String())
			dd.setErr(err)
		})

		dd.setErr(do(func(t T) { dd.setValue(t) }))
	}()
	return dd
}
