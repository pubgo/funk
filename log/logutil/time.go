package logutil

import (
	"strconv"
	"time"
)

const (
	timeFormatUnix      = ""
	timeFormatUnixMs    = "UNIXMS"
	timeFormatUnixMicro = "UNIXMICRO"
	timeFormatUnixNano  = "UNIXNANO"
)

// AppendTime formats the input time with the given format
// and appends the encoded string to the input byte slice.
func AppendTime(dst []byte, t time.Time, format string) []byte {
	switch format {
	case timeFormatUnix:
		return AppendInt64(dst, t.Unix())
	case timeFormatUnixMs:
		return AppendInt64(dst, t.UnixNano()/1000000)
	case timeFormatUnixMicro:
		return AppendInt64(dst, t.UnixNano()/1000)
	case timeFormatUnixNano:
		return AppendInt64(dst, t.UnixNano())
	}
	return append(t.AppendFormat(append(dst, '"'), format), '"')
}

// AppendDuration formats the input duration with the given unit & format
// and appends the encoded string to the input byte slice.
func AppendDuration(dst []byte, d time.Duration, unit time.Duration, useInt bool) []byte {
	if useInt {
		return strconv.AppendInt(dst, int64(d/unit), 10)
	}
	return AppendFloat64(dst, float64(d)/float64(unit))
}
