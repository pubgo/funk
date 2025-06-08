package anyhow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// === JSON Operations ===

// JsonMarshal wraps json.Marshal in a Result
func JsonMarshal[T any](v T) Result[[]byte] {
	return Try(func() ([]byte, error) {
		return json.Marshal(v)
	})
}

// JsonUnmarshal wraps json.Unmarshal in a Result
func JsonUnmarshal[T any](data []byte) Result[T] {
	return Try(func() (T, error) {
		var v T
		err := json.Unmarshal(data, &v)
		return v, err
	})
}

// JsonMarshalToString marshals to JSON and converts to string
func JsonMarshalToString[T any](v T) Result[string] {
	return MapTo(JsonMarshal(v), func(bytes []byte) string {
		return string(bytes)
	})
}

// JsonUnmarshalFromString unmarshals from JSON string
func JsonUnmarshalFromString[T any](s string) Result[T] {
	return JsonUnmarshal[T]([]byte(s))
}

// === File Operations ===

// ReadFile wraps os.ReadFile in a Result
func ReadFile(filename string) Result[[]byte] {
	return Try(func() ([]byte, error) {
		return os.ReadFile(filename)
	})
}

// WriteFile wraps os.WriteFile in a Result
func WriteFile(filename string, data []byte, perm os.FileMode) Result[struct{}] {
	return Try(func() (struct{}, error) {
		err := os.WriteFile(filename, data, perm)
		return struct{}{}, err
	})
}

// ReadTextFile reads a file and returns its content as string
func ReadTextFile(filename string) Result[string] {
	return MapTo(ReadFile(filename), func(bytes []byte) string {
		return string(bytes)
	})
}

// WriteTextFile writes string content to a file
func WriteTextFile(filename string, content string, perm os.FileMode) Result[struct{}] {
	return WriteFile(filename, []byte(content), perm)
}

// OpenFile wraps os.OpenFile in a Result
func OpenFile(name string, flag int, perm os.FileMode) Result[*os.File] {
	return Try(func() (*os.File, error) {
		return os.OpenFile(name, flag, perm)
	})
}

// === IO Operations ===

// ReadAll wraps io.ReadAll in a Result
func ReadAll(r io.Reader) Result[[]byte] {
	return Try(func() ([]byte, error) {
		return io.ReadAll(r)
	})
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) Result[int64] {
	return Try(func() (int64, error) {
		srcFile, err := os.Open(src)
		if err != nil {
			return 0, err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return 0, err
		}
		defer dstFile.Close()

		return io.Copy(dstFile, srcFile)
	})
}

// === String Conversions ===

// ParseInt wraps strconv.ParseInt in a Result
func ParseInt(s string, base int, bitSize int) Result[int64] {
	return Try(func() (int64, error) {
		return strconv.ParseInt(s, base, bitSize)
	})
}

// ParseFloat wraps strconv.ParseFloat in a Result
func ParseFloat(s string, bitSize int) Result[float64] {
	return Try(func() (float64, error) {
		return strconv.ParseFloat(s, bitSize)
	})
}

// ParseBool wraps strconv.ParseBool in a Result
func ParseBool(s string) Result[bool] {
	return Try(func() (bool, error) {
		return strconv.ParseBool(s)
	})
}

// Atoi wraps strconv.Atoi in a Result
func Atoi(s string) Result[int] {
	return Try(func() (int, error) {
		return strconv.Atoi(s)
	})
}

// === Formatting ===

// Sprintf safely formats a string, catching any panics
func Sprintf(format string, args ...interface{}) Result[string] {
	return CatchPanic(func() string {
		return fmt.Sprintf(format, args...)
	})
}

// === Slice Operations ===

// GetIndex safely gets an element at index from a slice
func GetIndex[T any](slice []T, index int) Result[T] {
	if index < 0 || index >= len(slice) {
		return Fail[T](fmt.Errorf("index %d out of bounds for slice of length %d", index, len(slice)))
	}
	return Ok(slice[index])
}

// FirstElement gets the first element of a slice
func FirstElement[T any](slice []T) Result[T] {
	return GetIndex(slice, 0)
}

// LastElement gets the last element of a slice
func LastElement[T any](slice []T) Result[T] {
	if len(slice) == 0 {
		return Fail[T](fmt.Errorf("slice is empty"))
	}
	return Ok(slice[len(slice)-1])
}

// === Map Operations ===

// GetKey safely gets a value from a map
func GetKey[K comparable, V any](m map[K]V, key K) Result[V] {
	if value, exists := m[key]; exists {
		return Ok(value)
	}
	return Fail[V](fmt.Errorf("key not found: %v", key))
}

// === Channel Operations ===

// ReceiveWithTimeout receives from a channel with timeout using select
func ReceiveWithTimeout[T any](ch <-chan T, timeoutCh <-chan struct{}) Result[T] {
	select {
	case value := <-ch:
		return Ok(value)
	case <-timeoutCh:
		return Fail[T](fmt.Errorf("timeout while receiving from channel"))
	}
}

// TryReceive attempts to receive from a channel without blocking
func TryReceive[T any](ch <-chan T) Result[T] {
	select {
	case value := <-ch:
		return Ok(value)
	default:
		return Fail[T](fmt.Errorf("no value available on channel"))
	}
}

// === Type Assertions ===

// TypeAssert safely performs type assertion
func TypeAssert[T any](value interface{}) Result[T] {
	if v, ok := value.(T); ok {
		return Ok(v)
	}
	return Fail[T](fmt.Errorf("type assertion failed: cannot convert %T to %T", value, *new(T)))
}

// === Validation Helpers ===

// NotEmpty validates that a string is not empty
func NotEmpty(s string) Result[string] {
	if s == "" {
		return Fail[string](fmt.Errorf("string is empty"))
	}
	return Ok(s)
}

// NotNil validates that a pointer is not nil
func NotNil[T any](ptr *T) Result[*T] {
	if ptr == nil {
		return Fail[*T](fmt.Errorf("pointer is nil"))
	}
	return Ok(ptr)
}

// InRange validates that a number is within a range
func InRange[T comparable](value, min, max T) Result[T] {
	// Note: This is a simplified version. In real use, you'd want proper numeric comparison
	return Ok(value) // Placeholder - would need proper implementation based on type constraints
}
