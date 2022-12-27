package log

import (
	"time"

	"github.com/pubgo/funk/stack"
)

var NameDelim = "."
var FieldCallerKey = "caller"
var FieldLoggerKey = "logger"
var FieldLevelKey = "level"
var FieldMsgKey = "msg"
var FieldErrorKey = "error"
var FieldErrorDetailKey = "error_detail"
var FieldStackTraceKey = "stacktrace"
var FieldTimestampKey = "ts"

// TimeField defines the time filed name in output.  It uses "time" in if empty.
var TimeField string

// TimeFormat specifies the time format in output. It uses RFC3339 with millisecond if empty.
// If set with `TimeFormatUnix/TimeFormatUnixMs/TimeFormatUnixWithMs`, timestamps are formated.
var TimeFormat string

var (
	// DefaultTimestamp is a Valuer that returns the current wallclock time,
	// respecting time zones, when bound.
	DefaultTimestamp = time.RFC3339Nano

	// DefaultCaller is a Valuer that returns the file and line where the Log
	// method was invoked. It can only be used with log.With.
	DefaultCaller = stack.Caller(0)
)

var DefaultFormatter string
