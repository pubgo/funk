package result

func IteratorOf[T any]() *Iterator[T] {
	return &Iterator[T]{v: make(chan Result[T])}
}

type Iterator[T any] struct {
	v    chan Result[T]
	done chan struct{}
}

func (cc *Iterator[T]) Done() {
	close(cc.v)
}

func (cc *Iterator[T]) Next() (Result[T], bool) {
	r, ok := <-cc.v
	return r, ok
}

func (cc *Iterator[T]) Range(fn func(r Result[T])) {
	for c := range cc.v {
		fn(c)
	}
}

func (cc *Iterator[T]) Chan() chan Result[T] {
	return cc.v
}
