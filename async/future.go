package async

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

func (f *Future[T]) close() {
	close(f.done)
}

func (f *Future[T]) setOK(v T) {
	defer f.close()
	f.v = v
}

func (f *Future[T]) setErr(err error) {
	defer f.close()
	f.e = err
}

func (f *Future[T]) Await() result.Result[T] {
	<-f.done
	return result.Wrap(f.v, f.e)
}
