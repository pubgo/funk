package async

import (
	"sync"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
)

func Promise[T any](fn func(resolve func(T), reject func(err error))) *result.Future[T] {
	assert.If(fn == nil, "[fn] is nil")

	var f = result.NewFuture[T]()
	go func() {
		defer recovery.Recovery(func(err errors.XError) {
			err.AddTag("fn", stack.CallerWithFunc(fn).String())
			f.Err(err)
		})

		fn(func(t T) { f.OK(t) }, func(err error) { f.Err(err) })
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
		dd.Chan() <- result.Err[T](try.Try(func() error {
			return do(func(t T) {
				dd.Chan() <- result.OK(t)
			})
		}))
	}()
	return dd
}
