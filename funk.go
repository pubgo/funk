package funk

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

//go:inline
func FromPtr[T any](v *T) (r T) {
	if v == nil {
		return
	}

	return *v
}

//go:inline
func ToPtr[T any](v T) *T {
	return &v
}

func Last[T any](args []T) (t T) {
	if len(args) == 0 {
		return
	}

	return args[len(args)-1]
}

func Ternary[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}

func TernaryFn[T any](ok bool, a, b func() T) T {
	if ok {
		return a()
	}
	return b()
}

func Map[T, V any](data []T, handle func(i int, d T) V) []V {
	vv := make([]V, 0, len(data))
	for i := range data {
		vv = append(vv, handle(i, data[i]))
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

// Filter iterates over `set` and gets the values that match `criteria`.
//
// Filter will return a new allocated slice.
func Filter[T any](set []T, checkTrue func(T) bool) []T {
	r := make([]T, 0)
	for i := range set {
		if !checkTrue(set[i]) {
			continue
		}
		r = append(r, set[i])
	}
	return r
}

// Delete the first occurrence of a type from a set.
func Delete[T comparable](set []T, value T) []T {
	for i := 0; i < len(set); i++ {
		if set[i] == value {
			set = append(set[:i], set[i+1:]...)
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

	v := reflect.ValueOf(err)
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Slice, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func Init(fn func()) error {
	fn()
	return nil
}

func DoFunc[T any](fn func() T) T {
	return fn()
}
