package clone

import (
	"reflect"
	"unsafe"

	"github.com/huandu/go-clone"
)

type Func = clone.Func
type Allocator = clone.Allocator
type AllocatorMethods = clone.AllocatorMethods

func Clone[T any](t T) T {
	return clone.Clone(t).(T)
}

func Slowly[T any](t T) T {
	return clone.Slowly(t).(T)
}

func Wrap[T any](t T) T {
	return clone.Wrap(t).(T)
}

func Unwrap[T any](t T) T {
	return clone.Unwrap(t).(T)
}

func Undo[T any](t T) {
	clone.Undo(t)
}

func MarkAsOpaquePointer(t reflect.Type) {
	clone.MarkAsOpaquePointer(t)
}

func MarkAsScalar(t reflect.Type) {
	clone.MarkAsScalar(t)
}

func SetCustomFunc(t reflect.Type, fn Func) {
	clone.SetCustomFunc(t, fn)
}

func FromHeap() *Allocator {
	return clone.FromHeap()
}

func NewAllocator(pool unsafe.Pointer, methods *AllocatorMethods) (allocator *Allocator) {
	return clone.NewAllocator(pool, methods)
}
