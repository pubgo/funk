package bytes

import "unsafe"

// ToString converts an array of bytes into a string without allocating.
// The byte slice passed to this function is not to be used after this call as it's unsafe; you have been warned.
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
