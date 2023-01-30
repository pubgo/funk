package async

import (
	"sync"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
)

func Promise[T any](fn func(resolve func(T), reject func(err error))) *Future[T] {
	assert.If(fn == nil, "[fn] is nil")

	var f = newFuture[T]()
	go func() {
		defer recovery.Recovery(func(err errors.XError) {
			err.AddTag("fn", stack.CallerWithFunc(fn).String())
			f.setErr(err)
		})

		fn(func(t T) { f.setOK(t) }, func(err error) { f.setErr(err) })
	}()
	return f
}

func Group[T any](do func(async func(func() (T, error))) error) *result.Iterator[T] {
	assert.If(do == nil, "[Async] [fn] is nil")

	var rr = result.IteratorOf[T]()
	go func() {
		var wg sync.WaitGroup
		defer rr.Done()
		defer wg.Wait()
		rr.Chan() <- result.Err[T](do(func(f func() (T, error)) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rr.Chan() <- result.Wrap(f())
			}()
		}))
	}()

	return rr
}

func Yield[T any](do func(yield func(T)) error) *result.Iterator[T] {
	var dd = result.IteratorOf[T]()
	go func() {
		defer dd.Done()
		defer recovery.Recovery(func(err errors.XError) {
			err.AddTag("fn_stack", stack.CallerWithFunc(do).String())
			dd.Chan() <- result.Err[T](err)
		})

		var err = do(func(t T) { dd.Chan() <- result.OK(t) })
		if err != nil {
			dd.Chan() <- result.Err[T](err)
		}
	}()
	return dd
}
