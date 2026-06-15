package async

import (
	"sync"

	"github.com/pubgo/funk/v2/result"
)

func iteratorOf[T any]() *Iterator[T] {
	return &Iterator[T]{v: make(chan T)}
}

type Iterator[T any] struct {
	v   chan T
	mu  sync.Mutex
	err error
}

func (cc *Iterator[T]) setDone() {
	close(cc.v)
}

func (cc *Iterator[T]) setErr(err error) {
	cc.mu.Lock()
	cc.err = err
	cc.mu.Unlock()
}

func (cc *Iterator[T]) setValue(v T) {
	cc.v <- v
}

func (cc *Iterator[T]) Next() (T, bool) {
	r, ok := <-cc.v
	return r, ok
}

func (cc *Iterator[T]) Await() result.Result[[]T] {
	ll := make([]T, 0, len(cc.v))
	for c := range cc.v {
		ll = append(ll, c)
	}

	cc.mu.Lock()
	err := cc.err
	cc.mu.Unlock()

	return result.Wrap(ll, err)
}
