package syncx

import (
	"github.com/pubgo/funk/result"
)

func newFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

type Future[T any] struct {
	v    T
	e    error
	done chan struct{}
}

func (f *Future[T]) success(v T) {
	defer close(f.done)
	f.v = v
}

func (f *Future[T]) failed(err error) {
	defer close(f.done)
	if err == nil {
		return
	}
	f.e = err
}

func (f *Future[T]) Await() result.Result[T] {
	<-f.done
	return result.New(f.v, f.e)
}

func (f *Future[T]) Unwrap() T {
	<-f.done
	if f.e == nil {
		return f.v
	}
	panic(f.e)
}

func (f *Future[T]) Value(check ...func(err error)) T {
	<-f.done
	if f.e != nil && len(check) > 0 && check[0] != nil {
		check[0](f.e)
	}
	return f.v
}
