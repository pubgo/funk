package log_config

import (
	"time"

	"github.com/pubgo/funk/stack"
)

// Defaults for Options.
const defaultTimestampFormat = "2006-01-02 15:04:05.000000"
const defaultMaxLogDepth = 16

var NameDelim = "."
var FieldCallerKey = "caller"
var FieldLoggerKey = "logger"
var FieldLevelKey = "level"
var FieldMsgKey = "msg"
var FieldErrorKey = "error"
var FieldErrorDetailKey = "error_detail"
var FieldStackTraceKey = "stacktrace"
var FieldTimestampKey = "ts"
var FieldTime = "time"

// "2006-01-02 15:04:05"
var TimeFormat = time.RFC3339

//atomic.AddUint64(&logNo, 1),

var (
	// DefaultTimestamp is a Valuer that returns the current wallclock time,
	// respecting time zones, when bound.
	DefaultTimestamp = time.RFC3339Nano
	defTimeFmt       = "2006-01-02 15:04:05"

	// DefaultCaller is a Valuer that returns the file and line where the Log
	// method was invoked. It can only be used with log.With.
	DefaultCaller = stack.Caller(0)
)

var DefaultFormatter string

// duration

//  Log as JSON instead of the default ASCII formatter.
//  log.SetFormatter(&log.JSONFormatter{})
//
//  // Output to stdout instead of the default stderr
//  // Can be any io.Writer, see below for File example
//  log.SetOutput(os.Stdout)
//
//  // Only log the warning severity or above.
//  log.SetLevel(log.WarnLevel)
