package reflectx

import (
	"reflect"
)

// Indirect returns the value pointed to by a pointer.
func Indirect(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		panic("[v] is invalid")
	}

	for {
		if v.Kind() != reflect.Ptr {
			return v
		}
		v = v.Elem()
	}
}

// New create a new value
func New(val any) reflect.Value {
	if val == nil {
		panic("[val] is nil")
	}

	return reflect.New(Indirect(reflect.ValueOf(val)).Type())
}

// FindFieldBy find field by handle
func FindFieldBy(v reflect.Value, handle func(field reflect.StructField) bool) reflect.Value {
	t := v.Type()
	for i := v.NumField() - 1; i >= 0; i-- {
		if handle(t.Field(i)) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}
