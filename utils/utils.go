package utils

import (
	"os"
	"strings"
	"time"

	"github.com/pubgo/funk/recovery"
)

func DotJoin(str ...string) string {
	return strings.Join(str, ".")
}

// DirExists function to check if directory exists?
func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		// path is a directory
		return true
	}
	return false
}

// FileExists function to check if file exists?
func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

func FirstFnNotEmpty(fx ...func() string) string {
	for i := range fx {
		if s := fx[i](); s != "" {
			return s
		}
	}
	return ""
}

func FirstNotEmpty(strL ...string) string {
	for i := range strL {
		if s := strL[i]; s != "" {
			return s
		}
	}
	return ""
}

func IfEmpty(str string, fx func()) {
	if str == "" {
		fx()
	}
}

func Cost(fn func()) (dur time.Duration, err error) {
	defer func(t time.Time) { dur = time.Since(t) }(time.Now())
	defer recovery.Err(&err)
	fn()
	return
}
