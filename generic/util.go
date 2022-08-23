package generic

import (
	"reflect"

	_ "golang.org/x/exp/constraints"
)

// Equals wraps the '==' operator for comparable types.
func Equals[T comparable](a, b T) bool {
	return a == b
}

func Nil[T any]() (t T) {
	return
}

func DePtr[T any](v *T) (r T) {
	if v == nil || reflect.ValueOf(*v).IsNil() {
		return
	}
	return *v
}

func Ptr[T any](v T) (r *T) {
	return &v
}

func Last[T any](args []T) (t T) {
	if len(args) == 0 {
		return
	}

	return args[len(args)-1]
}

func Ternary[T any](ok bool, a T, b T) T {
	if ok {
		return a
	}
	return b
}

func TernaryFn[T any](ok bool, a func() T, b func() T) T {
	if ok {
		return a()
	}
	return b()
}

func Map[T any](data []T, handle func(arg T) T) []T {
	for i := range data {
		data[i] = handle(data[i])
	}
	return data
}
