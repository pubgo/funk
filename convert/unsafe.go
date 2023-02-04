package convert

import "unsafe"

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func S2B(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// BtoS converts byte slice to string
func BtoS(b []byte) (s string) {
	B := (*Slice)(unsafe.Pointer(&b))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = B.Data
	S.Len = B.Len
	return
}

// StoB converts string to byte slice
func StoB(s string) (b []byte) {
	B := (*Slice)(unsafe.Pointer(&b))
	S := (*String)(unsafe.Pointer(&s))
	B.Data = S.Data
	B.Len = S.Len
	B.Cap = B.Len
	return
}
