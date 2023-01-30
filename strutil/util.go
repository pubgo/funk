package strutil

import (
	"reflect"
	"unsafe"
)

// ToBytes converts an existing string into an []byte without allocating.
// The string passed to this functions is not to be used again after this call as it's unsafe; you have been warned.
func ToBytes(s string) (b []byte) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHdr.Data = strHdr.Data
	sliceHdr.Cap = strHdr.Len
	sliceHdr.Len = strHdr.Len
	return
}

func FirstFnNotEmpty(ff ...func() string) string {
	for i := range ff {
		var v = ff[i]()
		if v != "" {
			return v
		}
	}
	return ""
}
