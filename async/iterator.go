package async

import "github.com/pubgo/funk/result"

func iteratorOf[T any]() *Iterator[T] {
	return &Iterator[T]{v: make(chan result.Result[T])}
}

type Iterator[T any] struct {
	v    chan result.Result[T]
	done chan struct{}
	err  error
}

func (cc *Iterator[T]) Err() error {
	return cc.err
}

func (cc *Iterator[T]) setDone() {
	close(cc.v)
}

func (cc *Iterator[T]) setErr(err error) {
	cc.err = err
}

func (cc *Iterator[T]) setValue(v T) {
	cc.v <- result.OK(v)
}

func (cc *Iterator[T]) Next() (result.Result[T], bool) {
	r, ok := <-cc.v
	return r, ok
}

func (cc *Iterator[T]) Range(fn func(r result.Result[T])) {
	for c := range cc.v {
		fn(c)
	}
}

func (cc *Iterator[T]) Chan() <-chan result.Result[T] {
	return cc.v
}
