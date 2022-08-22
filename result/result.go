package result

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pubgo/funk/generic"
)

func OK[T any](v T) Result[T] {
	return Result[T]{v: &v}
}

func Err[T any](err Error) Result[T] {
	return Result[T]{e: err}
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, e: WithErr(err)}
}

type Result[T any] struct {
	v *T
	e Error
}

func (r Result[T]) WithErr(err Error) Result[T] {
	r.e = err
	return r
}

func (r Result[T]) WithVal(t T) Result[T] {
	r.v = &t
	return r
}

func (r Result[T]) Err(check ...func(t T)) Error {
	if len(check) > 0 && check[0] != nil && !r.IsErr() {
		check[0](generic.DePtr(r.v))
		return Error{}
	}
	return r.e
}

func (r Result[T]) Value(check ...func(err Error) T) T {
	if len(check) > 0 && check[0] != nil && !r.e.IsNil() {
		return check[0](r.e)
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) IsErr() bool {
	return !r.e.IsNil()
}

func (r Result[T]) IsNil() bool {
	return r.v == nil || reflect.ValueOf(*r.v).IsNil()
}

func (r Result[T]) IsNone() bool {
	return !r.IsErr() && r.IsNil()
}

func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsErr() {
		return r
	}

	r.v = generic.Ptr(f(generic.DePtr(r.v)))
	return r
}

func (r Result[T]) Expect(msg string, args ...interface{}) T {
	if r.IsErr() {
		panic(r.e.WrapF(msg, args...))
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) UnwrapOr(v T) T {
	if r.IsErr() {
		return v
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) UnwrapOrElse(check ...func(err Error) T) T {
	if len(check) > 0 && check[0] != nil && r.IsErr() {
		return check[0](r.e)
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) Unwrap(check ...func(err Error) T) T {
	if r.IsErr() {
		if len(check) > 0 && check[0] != nil {
			return check[0](r.e)
		} else {
			panic(r.e)
		}
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%#v", r.v)
	}
	return fmt.Sprintf("err_msg=%q err_detail=%#v", r.e.Unwrap(), r.e.Unwrap())
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	var d = data[T]{Body: generic.DePtr(r.v)}
	if r.IsErr() {
		d.ErrMsg = r.e.Error()
		d.ErrDetail = fmt.Sprintf("%#v", r.e.Unwrap())
	}
	return json.Marshal(d)
}

type Chan[T any] chan Result[T]

func (cc Chan[T]) Unwrap(check ...func(err Error) []T) []T {
	return cc.ToResult().Unwrap(check...)
}

func (cc Chan[T]) ToList() List[T] {
	var rr []Result[T]
	for r := range cc {
		if r.IsNone() {
			continue
		}

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

		if c.IsNil() {
			continue
		}

		rl = append(rl, c.Value())
	}
	return OK(rl)
}

func (cc Chan[T]) Range(fn func(r Result[T])) {
	for c := range cc {
		if c.IsNone() {
			continue
		}

		fn(c)
	}
}

type List[T any] []Result[T]

func (rr List[T]) Unwrap(check ...func(err Error) []T) []T {
	return rr.ToResult().Unwrap(check...)
}

func (rr List[T]) Map(h func(r Result[T]) Result[T]) List[T] {
	var ll = make(List[T], 0, len(rr))
	for i := range rr {
		if rr[i].IsNone() {
			continue
		}
		ll = append(ll, h(rr[i]))
	}
	return ll
}

func (rr List[T]) Filter(filter func(r *Result[T]) bool) List[T] {
	var ll = make(List[T], 0, len(rr))
	for i := range rr {
		if rr[i].IsNone() {
			continue
		}

		if filter(&rr[i]) {
			ll = append(ll, rr[i])
		}
	}
	return ll
}

func (rr List[T]) ToResult() Result[[]T] {
	var rl = make([]T, 0, len(rr))
	for i := range rr {
		if rr[i].IsErr() {
			return Err[[]T](rr[i].Err())
		}
		if rr[i].IsNone() {
			continue
		}
		rl = append(rl, rr[i].Value())
	}
	return OK(rl)
}

func (rr List[T]) Range(fn func(r Result[T])) {
	for i := range rr {
		if rr[i].IsNone() {
			continue
		}
		fn(rr[i])
	}
}
