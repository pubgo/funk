package syncx

import (
	"sync"
	
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/internal/utils"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/xerr"
	"github.com/pubgo/funk/xtry"
)

func Promise[T any](fn func(resolve func(T), reject func(error))) *Future[T] {
	assert.If(fn == nil, "[fn] is nil")

	var f = newFuture[T]()
	go func() {
		defer recovery.Recovery(func(err xerr.XErr) {
			f.failed(err.WrapF("fn=%s", utils.CallerWithFunc(fn)))
		})
		fn(func(t T) { f.success(t) }, func(err error) {
			if err == nil {
				return
			}
			f.failed(err)
		})
	}()
	return f
}

func AsyncGroup[T any](do func(async func(func() result.Result[T])) error) result.Chan[T] {
	assert.If(do == nil, "[Async] [fn] is nil")

	var rr = make(chan result.Result[T])
	go func() {
		var wg sync.WaitGroup
		var err = xtry.Try(func() error {
			return do(func(f func() result.Result[T]) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					rr <- result.OK(xtry.TryCatch1(f, func(err xerr.XErr) {
						rr <- result.Err[T](err)
					}))
				}()
			})
		})
		if err != nil {
			rr <- result.Err[T](err)
		}
		wg.Wait()
		close(rr)
	}()

	return rr
}

func Wait[T any](val ...*Future[T]) result.List[T] {
	var valList = make([]result.Result[T], len(val))
	for i := range val {
		valList[i] = val[i].Await()
	}
	return valList
}

func Yield[T any](do func(yield func(T)) error) result.Chan[T] {
	var dd = make(chan result.Result[T])
	go func() {
		defer close(dd)
		err := xtry.Try(func() error { return do(func(t T) { dd <- result.OK(t) }) })
		if err != nil {
			dd <- result.Err[T](err)
		}
	}()
	return dd
}
