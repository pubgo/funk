package result

import "github.com/pubgo/funk/errors"

func ChanOf[T any](args chan T) *Chan[T] {
	var ret = &Chan[T]{v: make(chan T)}
	go func() {
		defer close(ret.v)
		for arg := range args {
			ret.v <- arg
		}
	}()
	return ret
}

type Chan[T any] struct {
	v chan T
	e error
}

func (cc *Chan[T]) SetErr(err error) *Chan[T] {
	cc.e = err
	return cc
}

func (cc *Chan[T]) Err() error {
	return cc.e
}

func (cc *Chan[T]) IsErr() bool {
	return errors.IsNil(cc.e)
}

func (cc *Chan[T]) Unwrap() chan T {
	return cc.v
}

func (cc *Chan[T]) Range(fn func(r T)) {
	for c := range cc.v {
		fn(c)
	}
}
