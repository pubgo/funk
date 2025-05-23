package log

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

var _ zerolog.LogObjectMarshaler = (*logLogObjectMarshaler)(nil)

type logLogObjectMarshaler struct {
	err error
}

func (l logLogObjectMarshaler) MarshalZerologObject(e *zerolog.Event) {
	if l.err == nil {
		return
	}

	if errJson, ok := l.err.(json.Marshaler); ok {
		errJsonBytes, _ := errJson.MarshalJSON()
		if len(errJsonBytes) > 0 {
			e.Str("msg", l.err.Error())
			e.RawJSON("detail", errJsonBytes)
			return
		}
	}

	e.Str("error", l.err.Error())
	e.Str("error_detail", fmt.Sprintf("%v", l.err))
}
