package result

import "github.com/pubgo/funk/errors"

func ListOf[T any](args ...T) List[T] {
	var ret = List[T]{v: make([]T, 0, len(args))}
	for i := range args {
		ret.v = append(ret.v, args[i])
	}
	return ret
}

type List[T any] struct {
	v []T
	e error
}

func (rr List[T]) Err() error {
	return rr.e
}

func (rr List[T]) IsErr() bool {
	return errors.IsNil(rr.e)
}

func (rr List[T]) Unwrap() []T {
	return rr.v
}

func (rr List[T]) ToResult() (r Result[[]T]) {
	if rr.IsErr() {
		return r.WithErr(rr.e)
	}
	return r.WithVal(rr.v)
}

func (rr List[T]) Range(fn func(r T)) {
	for i := range rr.v {
		fn(rr.v[i])
	}
}
