package async

import (
	"github.com/pubgo/funk/result"
)

func iteratorOf[T any]() *Iterator[T] {
	return &Iterator[T]{v: make(chan T)}
}

type Iterator[T any] struct {
	v   chan T
	err error
}

func (cc *Iterator[T]) setDone() {
	close(cc.v)
}

func (cc *Iterator[T]) setErr(err error) {
	cc.err = err
}

func (cc *Iterator[T]) setValue(v T) {
	cc.v <- v
}

func (cc *Iterator[T]) Next() (T, bool) {
	r, ok := <-cc.v
	return r, ok
}

func (cc *Iterator[T]) Await() result.Result[[]T] {
	var ll []T
	err := cc.err
	if err != nil {
		return result.Wrap(ll, err)
	}

	for c := range cc.v {
		ll = append(ll, c)
	}
	return result.Wrap(ll, cc.err)
}
