package log

import "time"

var NameDelim = "."
var FieldCallerKey = "caller"
var FieldLoggerKey = "logger"
var FieldLevelKey = "level"
var FieldMsgKey = "msg"
var FieldErrorKey = "error"
var FieldErrorDetailKey = "error_detail"
var FieldStackTraceKey = "stacktrace"
var FieldTimestampKey = "ts"
var (
	// DefaultTimestamp is a Valuer that returns the current wallclock time,
	// respecting time zones, when bound.
	DefaultTimestamp = TimestampFormat(time.Now, time.RFC3339Nano)

	// DefaultTimestampUTC is a Valuer that returns the current time in UTC
	// when bound.
	DefaultTimestampUTC = TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		time.RFC3339Nano,
	)

	// DefaultCaller is a Valuer that returns the file and line where the Log
	// method was invoked. It can only be used with log.With.
	DefaultCaller = Caller(3)
)
