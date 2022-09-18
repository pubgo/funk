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

func Promise[T any](fn func(resolve func(T), reject func(result.Error))) *Future[T] {
	assert.If(fn == nil, "[fn] is nil")

	var f = newFuture[T]()
	go func() {
		defer recovery.Recovery(func(err xerr.XErr) {
			f.failed(result.WithErr(err).WrapF("fn=%s", utils.CallerWithFunc(fn)))
		})

		fn(func(t T) { f.success(result.OK(t)) }, f.failed)
	}()
	return f
}

func AsyncGroup[T any](do func(async func(func() result.Result[T])) result.Error) result.Chan[T] {
	assert.If(do == nil, "[Async] [fn] is nil")

	var rr = make(chan result.Result[T])
	go func() {
		var wg sync.WaitGroup
		rr <- result.Err[T](xtry.TryErr(func() result.Error {
			return do(func(f func() result.Result[T]) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					rr <- xtry.TryVal(f)
				}()
			})
		}).Err())
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

func Yield[T any](do func(yield func(T)) result.Error) result.Chan[T] {
	var dd = make(chan result.Result[T])
	go func() {
		defer close(dd)
		err := xtry.TryErr(func() result.Error { return do(func(t T) { dd <- result.OK(t) }) })
		dd <- result.Err[T](err.Err())
	}()
	return dd
}
