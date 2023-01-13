package tracedebug

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/pubgo/funk/debug"
	"golang.org/x/net/trace"
)

func init() {
	debug.Get("/requests", adaptor.HTTPHandlerFunc(trace.Traces))
	debug.Get("/events", adaptor.HTTPHandlerFunc(trace.Events))
}
