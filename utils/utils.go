package utils

import (
	"time"

	"github.com/pubgo/funk/recovery"
)

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
