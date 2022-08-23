package syncx

import (
	"github.com/pubgo/funk/result"
)

func newFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

type Future[T any] struct {
	v    result.Result[T]
	done chan struct{}
}

func (f *Future[T]) success(v result.Result[T]) {
	defer close(f.done)
	f.v = v
}

func (f *Future[T]) failed(err result.Error) {
	defer close(f.done)
	if err.IsNil() {
		return
	}
	f.v = result.Err[T](err)
}

func (f *Future[T]) Await() result.Result[T] {
	<-f.done
	return f.v
}

func (f *Future[T]) Unwrap(check ...func(err result.Error) T) T {
	<-f.done
	return f.v.Unwrap(check...)
}
