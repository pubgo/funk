package result

import (
	"fmt"

	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/logfields"
)

type Event struct {
	*log.Event
}

func (e Event) Msg(msg string) {
	e.Str(logfields.Msg, msg)
}

func (e Event) MsgFunc(createMsg func() string) {
	e.Str(logfields.Msg, createMsg())
}

func (e Event) Msgf(format string, v ...any) {
	e.Str(logfields.Msg, fmt.Sprintf(format, v...))
}
