package logx

import (
	logkit "github.com/go-kit/log"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/go-logfmt/logfmt"

	"github.com/go-logr/logr"
	"github.com/iand/logfmtr"
	"github.com/pubgo/funk/logger"
)

var logWriter logger.Logger
var gl = logr.Discard()
var TimestampFormat = time.RFC3339
var NameDelim = "."
var DefaultCallerSkip = 2

// The global verbosity level.
var gv int32 = 2

func init() {
	opts := logfmtr.DefaultOptions()
	opts.Writer = os.Stderr
	opts.Humanize = true
	opts.Colorize = true
	opts.CallerSkip = DefaultCallerSkip
	opts.AddCaller = true

	logkit.Caller()

	var log = logfmtr.NewWithOptions(opts)
	gl = logr.New(&sink{})
}

func SetLogWriter(w logger.Logger) {
	logWriter = w
}

// SetVerbosity sets the global log level.
//
//	Only loggers with a V level less than or equal to this value will be enabled.
func SetVerbosity(v int) int {
	old := atomic.SwapInt32(&gv, int32(v))
	return int(old)
}
