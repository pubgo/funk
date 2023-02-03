package generic

import (
	"reflect"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func AppendOf[T any](v T, vv ...T) []T {
	return append(append(make([]T, 0, len(vv)+1), v), vv...)
}

func ListOf[T any](args ...T) []T {
	return args
}

func Zero[T any]() (ret T) {
	return
}

// Equals wraps the '==' operator for comparable types.
func Equals[T comparable](a, b T) bool {
	return a == b
}

func Nil[T any]() (t *T) {
	return
}

func DePtr[T any](v *T) (r T) {
	if v == nil || reflect.ValueOf(*v).IsNil() {
		return
	}
	return *v
}

//go:inline
func Ptr[T any](v T) *T {
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

func Map[T any, V any](data []T, handle func(i int) V) []V {
	var vv = make([]V, 0, len(data))
	for i := range data {
		vv = append(vv, handle(i))
	}
	return vv
}

// Contains returns whether `vs` contains the element `e` by comparing vs[i] == e.
func Contains[T comparable](vs []T, e T) bool {
	for _, v := range vs {
		if v == e {
			return true
		}
	}

	return false
}

// ExtractFrom extracts a nested object of type E from type T.
//
// This function is useful if we have a set of type `T` nad we want to
// extract the type E from any T.
func ExtractFrom[T, E any](set []T, fn func(T) E) []E {
	r := make([]E, len(set))
	for i := range set {
		r[i] = fn(set[i])
	}

	return r
}

// Filter iterates over `set` and gets the values that match `criteria`.
//
// Filter will return a new allocated slice.
func Filter[T any](set []T, criteria func(T) bool) []T {
	r := make([]T, 0)
	for i := range set {
		if criteria(set[i]) {
			r = append(r, set[i])
		}
	}

	return r
}

// FilterInPlace filters the contents of `set` using `criteria`.
//
// FilterInPlace returns `set`.
func FilterInPlace[T any](set []T, criteria func(T) bool) []T {
	for i := 0; i < len(set); i++ {
		if !criteria(set[i]) {
			set = append(set[:i], set[i+1:]...)
			i--
		}
	}

	return set
}

// Delete the first occurrence of a type from a set.
func Delete[T comparable](set []T, value T) []T {
	for i := 0; i < len(set); i++ {
		if set[i] == value {
			set = append(set[:i], set[i:]...)
			break
		}
	}

	return set
}

// DeleteAll occurrences from a set.
func DeleteAll[T comparable](set []T, value T) []T {
	for i := 0; i < len(set); i++ {
		if set[i] == value {
			set = append(set[:i], set[i:]...)
			i--
		}
	}

	return set
}

// Max returns the max of the 2 passed values.
func Max[T constraints.Ordered](a, b T) (r T) {
	if a < b {
		r = b
	} else {
		r = a
	}

	return
}

// Min returns the min of the 2 passed values.
func Min[T constraints.Ordered](a, b T) (r T) {
	if a < b {
		r = a
	} else {
		r = b
	}

	return
}

// isNilValue copy from <github.com/rs/zerolog.isNilValue>
func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

func IsNil(err interface{}) bool {
	if err == nil {
		return true
	}

	if isNilValue(err) {
		return true
	}

	var v = reflect.ValueOf(err)
	if !v.IsValid() {
		return true
	}

	return v.IsZero()
}
