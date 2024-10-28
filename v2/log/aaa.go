package log

import (
	"github.com/pubgo/funk/v2/log/log_internal"
)

type Logger = log_internal.Logger
type Event = log_internal.Event
type Map = log_internal.Map
type EnableChecker = log_internal.EnableChecker
type StdLogger = log_internal.StdLogger

var CreateEventCtx = log_internal.CreateEventCtx
var UpdateEventCtx = log_internal.UpdateEventCtx
var GetEventFromCtx = log_internal.GetEventFromCtx
var WithDisabled = log_internal.WithDisabled
var NewStd = log_internal.NewStd
