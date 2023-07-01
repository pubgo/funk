package convert

import (
	"unsafe"
)

// CopyBytes copies a slice to make it immutable
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// Copy creates an identical copy of x
func Copy(x []byte) []byte { return x[:len(x):len(x)] }

// CopyString copies a string to make it immutable
func CopyString(s string) string {
	return string(S2B(s))
}

// String internals from reflect
type String struct {
	Data unsafe.Pointer
	Len  int
}

// Slice internals from reflect
type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// BtoU4 converts byte slice to integer slice
func BtoU4(b []byte) (i []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 2
	I.Cap = I.Len
	return
}

// U4toB converts integer slice to byte slice
func U4toB(i []uint32) (b []byte) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 2
	B.Cap = B.Len
	return
}

// U4toU8 converts uint32 slice to uint64 slice
func U4toU8(i []uint32) (k []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	K := (*Slice)(unsafe.Pointer(&k))
	K.Data = I.Data
	K.Len = I.Len >> 1
	K.Cap = K.Len
	return
}

// U8toU4 converts uint64 slice to uint32 slice
func U8toU4(i []uint64) (k []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	K := (*Slice)(unsafe.Pointer(&k))
	K.Data = I.Data
	K.Len = I.Len << 1
	K.Cap = K.Len
	return
}

// BtoU8 converts byte slice to integer slice
func BtoU8(b []byte) (i []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 3
	I.Cap = I.Len
	return
}

// U8toB converts integer slice to byte slice
func U8toB(i []uint64) (b []byte) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 3
	B.Cap = B.Len
	return
}

// StoU4 converts string to integer slice
func StoU4(s string) (i []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 2
	I.Cap = I.Len
	return
}

// U4toS converts integer slice to string
func U4toS(i []uint32) (s string) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 2
	return
}

// StoU8 converts string to integer slice
func StoU8(s string) (i []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 3
	I.Cap = I.Len
	return
}

// U8toS converts integer slice to string
func U8toS(i []uint64) (s string) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 3
	return
}
