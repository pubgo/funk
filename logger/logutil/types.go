package logutil

import (
	"math"
	"strconv"
)

// AppendBool converts the input bool to a string and
// appends the encoded string to the input byte slice.
func AppendBool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

// AppendInt converts the input int to a string and
// appends the encoded string to the input byte slice.
func AppendInt(dst []byte, val int) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInt8 converts the input []int8 to a string and
// appends the encoded string to the input byte slice.
func AppendInt8(dst []byte, val int8) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInt16 converts the input int16 to a string and
// appends the encoded string to the input byte slice.
func AppendInt16(dst []byte, val int16) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInt32 converts the input int32 to a string and
// appends the encoded string to the input byte slice.
func AppendInt32(dst []byte, val int32) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

// AppendInt64 converts the input int64 to a string and
// appends the encoded string to the input byte slice.
func AppendInt64(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, val, 10)
}

// AppendUint converts the input uint to a string and
// appends the encoded string to the input byte slice.
func AppendUint(dst []byte, val uint) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUint8 converts the input uint8 to a string and
// appends the encoded string to the input byte slice.
func AppendUint8(dst []byte, val uint8) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUint16 converts the input uint16 to a string and
// appends the encoded string to the input byte slice.
func AppendUint16(dst []byte, val uint16) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUint32 converts the input uint32 to a string and
// appends the encoded string to the input byte slice.
func AppendUint32(dst []byte, val uint32) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

// AppendUint64 converts the input uint64 to a string and
// appends the encoded string to the input byte slice.
func AppendUint64(dst []byte, val uint64) []byte {
	return strconv.AppendUint(dst, val, 10)
}

func appendFloat(dst []byte, val float64, bitSize int) []byte {
	// JSON does not permit NaN or Infinity. A typical JSON encoder would fail
	// with an error, but a logging library wants the data to get through so we
	// make a tradeoff and store those types as string.
	switch {
	case math.IsNaN(val):
		return append(dst, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(dst, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(dst, `"-Inf"`...)
	}
	return strconv.AppendFloat(dst, val, 'f', -1, bitSize)
}

// AppendFloat32 converts the input float32 to a string and
// appends the encoded string to the input byte slice.
func AppendFloat32(dst []byte, val float32) []byte {
	return appendFloat(dst, float64(val), 32)
}

// AppendFloat64 converts the input float64 to a string and
// appends the encoded string to the input byte slice.
func AppendFloat64(dst []byte, val float64) []byte {
	return appendFloat(dst, val, 64)
}
