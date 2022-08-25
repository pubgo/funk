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
	err.Must()
	return Result[T]{e: err}
}

func Wrap[T any](v T, err error) Result[T] {
	return Result[T]{v: &v, e: WithErr(err)}
}

type Result[T any] struct {
	v *T
	e Error
}

func (r Result[T]) Err(check ...func(t T)) Error {
	if len(check) > 0 && check[0] != nil && !r.IsErr() {
		check[0](generic.DePtr(r.v))
		return Error{}
	}
	return r.e
}

func (r Result[T]) IsErr() bool {
	return r.e.IsErr()
}

func (r Result[T]) IsNil() bool {
	return r.v == nil || reflect.ValueOf(*r.v).IsNil()
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

func (r Result[T]) OrElse(v T) T {
	if r.IsErr() {
		return v
	}
	return generic.DePtr(r.v)
}

func (r Result[T]) Unwrap(check ...func(err Error) T) T {
	if !r.IsErr() {
		return generic.DePtr(r.v)
	}

	if len(check) > 0 && check[0] != nil {
		return check[0](r.e)
	} else {
		panic(r.e)
	}
}

func (r Result[T]) String() string {
	if !r.IsErr() {
		return fmt.Sprintf("%v", r.v)
	}
	return fmt.Sprintf("err_msg=%q err_detail=%#v", r.e.Unwrap(), r.e.Unwrap())
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	var d = data[T]{Body: generic.DePtr(r.v)}
	if r.IsErr() {
		d.ErrMsg = r.e.Err().Error()
		d.ErrDetail = fmt.Sprintf("%#v", r.e.Unwrap())
	}
	return json.Marshal(d)
}

func ChanOf[T any](args chan T) Chan[T] {
	var ret = make(Chan[T])
	go func() {
		defer close(ret)
		for arg := range args {
			ret <- OK(arg)
		}
	}()
	return ret
}

type Chan[T any] chan Result[T]

func (cc Chan[T]) Unwrap(check ...func(err Error) []T) []T {
	return cc.ToResult().Unwrap(check...)
}

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

		if c.IsNil() {
			continue
		}

		rl = append(rl, c.Unwrap())
	}
	return OK(rl)
}

func (cc Chan[T]) Range(fn func(r Result[T])) {
	for c := range cc {
		fn(c)
	}
}

func ListOf[T any](args ...T) List[T] {
	var ret = make([]Result[T], 0, len(args))
	for i := range args {
		ret = append(ret, OK(args[i]))
	}
	return ret
}

type List[T any] []Result[T]

func (rr List[T]) Unwrap(check ...func(err Error) []T) []T {
	return rr.ToResult().Unwrap(check...)
}

func (rr List[T]) Map(h func(r Result[T]) Result[T]) List[T] {
	var ll = make(List[T], 0, len(rr))
	for i := range rr {
		ll = append(ll, h(rr[i]))
	}
	return ll
}

func (rr List[T]) Filter(filter func(r *Result[T]) bool) List[T] {
	var ll = make(List[T], 0, len(rr))
	for i := range rr {
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
		rl = append(rl, rr[i].Unwrap())
	}
	return OK(rl)
}

func (rr List[T]) Range(fn func(r Result[T])) {
	for i := range rr {
		fn(rr[i])
	}
}
