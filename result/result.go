package result

import (
	"encoding/json"
	"fmt"
)

func Unwrap[T any](r Result[T]) (T, error) {
	return r.v, r.e
}

func OK[T any](v T) Result[T] {
	return Result[T]{v: v}
}

func Err[T any](err error) Result[T] {
	return Result[T]{e: err}
}

func New[T any](v T, err error) Result[T] {
	return Result[T]{v: v, e: err}
}

type Result[T any] struct {
	v T
	e error
}

func (r Result[T]) WithErr(err error) Result[T] {
	r.e = err
	return r
}

func (r Result[T]) WithVal(t T) Result[T] {
	r.v = t
	return r
}

func (r Result[T]) Err() error {
	return r.e
}

func (r Result[T]) IsErr() bool {
	return r.e != nil
}

func (r Result[T]) Must() T {
	if r.IsErr() {
		panic(r.e)
	}
	return r.v
}

func (r Result[T]) String() string {
	if r.e == nil {
		return fmt.Sprintf("v: %#v", r.v)
	}
	return fmt.Sprintf("err_msg=%s err_detail=%#v", r.e.Error(), r.e)
}

func (r Result[T]) Value(check ...func(err error)) T {
	if len(check) > 0 && check[0] != nil && r.e != nil {
		check[0](r.e)
	}
	return r.v
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	var err = EmptyErr()
	if r.e != nil {
		err = r.e
	}

	return json.Marshal(data[T]{
		Body:      r.v,
		ErrMsg:    err.Error(),
		ErrDetail: fmt.Sprintf("%#v", err),
	})
}

type Chan[T any] chan Result[T]

func (cc Chan[T]) ToList() List[T] {
	var rr []Result[T]
	for r := range cc {
		rr = append(rr, r)
	}
	return rr
}

func (cc Chan[T]) ToResult() Result[[]T] {
	var rl = make([]T, 0, len(cc))
	for c := range cc {
		if c.IsErr() {
			return Err[[]T](c.Err())
		}
		rl = append(rl, c.Value())
	}
	return OK(rl)
}

func (cc Chan[T]) Range(fn func(r Result[T])) {
	for c := range cc {
		fn(c)
	}
}

type List[T any] []Result[T]

func (rr List[T]) ToResult() Result[[]T] {
	var rl = make([]T, 0, len(rr))
	for i := range rr {
		if rr[i].IsErr() {
			return Err[[]T](rr[i].Err())
		}
		rl = append(rl, rr[i].Value())
	}
	return OK(rl)
}

func (rr List[T]) Range(fn func(r Result[T])) {
	for i := range rr {
		fn(rr[i])
	}
}
