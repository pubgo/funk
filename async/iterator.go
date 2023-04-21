package async

import "github.com/pubgo/funk/result"

func iteratorOf[T any]() *Iterator[T] {
	return &Iterator[T]{v: make(chan T)}
}

type Iterator[T any] struct {
	v    chan T
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
	cc.v <- v
}

func (cc *Iterator[T]) Next() (T, bool) {
	r, ok := <-cc.v
	return r, ok
}

func (cc *Iterator[T]) Range(fn func(r T)) error {
	for c := range cc.v {
		if cc.err != nil {
			return cc.err
		}

		fn(c)
	}
	return nil
}

func (cc *Iterator[T]) ToList() result.Result[[]T] {
	var ret result.Result[[]T]
	if cc.err != nil {
		return ret.WithErr(cc.err)
	}

	var ll []T
	for c := range cc.v {
		c1 := c
		ll = append(ll, c1)
	}
	return ret
}
