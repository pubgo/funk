package log_internal

import (
	"context"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
)

const ModuleName = "module"

var (
	globalLevel      Level
	logEnableChecker        = func(ctx context.Context, lvl Level, nameOrMessage string, fields Map) bool { return true }
	stdZeroLog       Logger = nil
	stdLog                  = New(nil)
)

func GetLogger() Logger { return stdLog }

// SetLogger set global log
func SetLogger(log EventLogger) {
	assert.If(log == nil, "[log] should not be nil")

	log = generic.Ptr(log.Hook(logGlobalHook))

	stdZeroLog = log
}

func SetEnableChecker(checker EnableChecker) {
	if checker == nil {
		return
	}

	logEnableChecker = checker
}
