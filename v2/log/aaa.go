package log

import (
	"github.com/pubgo/funk/v2/log/log_internal"
)

type Logger = log_internal.Logger
type Event = log_internal.Event
type Map = log_internal.Map
type EnableChecker = log_internal.EnableChecker
type StdLogger = log_internal.StdLogger

var CreateCtx = log_internal.CreateCtx
var UpdateCtx = log_internal.UpdateCtx
var GetFromCtx = log_internal.GetFromCtx
var WithDisabled = log_internal.WithDisabled
var NewStd = log_internal.NewStd
