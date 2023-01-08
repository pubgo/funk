package result

func NewFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

type Future[T any] struct {
	v    Result[T]
	done chan struct{}
}

func (f *Future[T]) OK(v T) {
	defer close(f.done)
	f.v = f.v.WithVal(v)
}

func (f *Future[T]) Result(r Result[T]) {
	defer close(f.done)
	f.v = r
}

func (f *Future[T]) Err(err error) {
	defer close(f.done)
	f.v = f.v.WithErr(err)
}

func (f *Future[T]) Await() Result[T] {
	<-f.done
	return f.v
}
