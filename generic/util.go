package generic

import (
	"github.com/pubgo/funk"
	"golang.org/x/exp/constraints"
)

func AppendOf[T any](v T, vv ...T) []T {
	return funk.AppendOf(v, vv...)
}

func ListOf[T any](args ...T) []T {
	return funk.ListOf(args...)
}

func Zero[T any]() (ret T) {
	return funk.Zero[T]()
}

func Equals[T comparable](a, b T) bool {
	return funk.Equals(a, b)
}

func Nil[T any]() (t *T) {
	return funk.Nil[T]()
}

// DePtr
// Deprecated: use FromPtr
func DePtr[T any](v *T) (r T) {
	return funk.FromPtr(v)
}

func FromPtr[T any](v *T) (r T) {
	return funk.FromPtr(v)
}

//go:inline
func Ptr[T any](v T) *T {
	return funk.ToPtr(v)
}

func Last[T any](args []T) (t T) {
	return funk.Last(args)
}

func Ternary[T any](ok bool, a, b T) T {
	return funk.Ternary(ok, a, b)
}

func TernaryFn[T any](ok bool, a, b func() T) T {
	return funk.TernaryFn(ok, a, b)
}

func Map[T, V any](data []T, handle func(i int) V) []V {
	vv := make([]V, 0, len(data))
	for i := range data {
		vv = append(vv, handle(i))
	}
	return vv
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

func IsNil(err interface{}) bool {
	return funk.IsNil(err)
}

func Init(fn func()) any {
	fn()
	return nil
}
