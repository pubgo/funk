package async

import "github.com/pubgo/funk/result"

func newFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

type Future[T any] struct {
	v    result.Result[T]
	done chan struct{}
}

func (f *Future[T]) close() {
	close(f.done)
}

func (f *Future[T]) setOK(v T) {
	defer f.close()
	f.v = f.v.WithVal(v)
}

func (f *Future[T]) setErr(err error) {
	defer f.close()
	f.v = f.v.WithErr(err)
}

func (f *Future[T]) Await() result.Result[T] {
	<-f.done
	return f.v
}
