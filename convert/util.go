package convert

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Map[K comparable, S, D any](src map[K]S, convert func(S) D) map[K]D {
	if src == nil {
		return nil
	}
	dst := make(map[K]D, len(src))
	for k, s := range src {
		dst[k] = convert(s)
	}
	return dst
}

func MapL[A comparable, B any](src []A, convert func(A) B) []B {
	if src == nil {
		return nil
	}

	dst := make([]B, len(src))
	for i := range src {
		dst[i] = convert(src[i])
	}
	return dst
}

const (
	uByte = 1 << (10 * iota)
	uKilobyte
	uMegabyte
	uGigabyte
	uTerabyte
	uPetabyte
	uExabyte
)

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)
	switch {
	case bytes >= uExabyte:
		unit = "EB"
		value /= uExabyte
	case bytes >= uPetabyte:
		unit = "PB"
		value /= uPetabyte
	case bytes >= uTerabyte:
		unit = "TB"
		value /= uTerabyte
	case bytes >= uGigabyte:
		unit = "GB"
		value /= uGigabyte
	case bytes >= uMegabyte:
		unit = "MB"
		value /= uMegabyte
	case bytes >= uKilobyte:
		unit = "KB"
		value /= uKilobyte
	case bytes >= uByte:
		unit = "B"
	default:
		return "0B"
	}
	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// ToString Change arg to string
func ToString(arg any, timeFormat ...string) string {
	tmp := reflect.Indirect(reflect.ValueOf(arg)).Interface()
	switch v := tmp.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeFormat) > 0 {
			return v.Format(timeFormat[0])
		}
		return v.Format("2006-01-02 15:04:05")
	case reflect.Value:
		return ToString(v.Interface(), timeFormat...)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
