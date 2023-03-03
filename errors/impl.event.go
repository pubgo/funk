package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/rs/zerolog"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

var _ ErrEvent = (*errEventImpl)(nil)
var _ fmt.Formatter = (*errEventImpl)(nil)

type errEventImpl struct {
	err    error
	caller *stack.Frame
	evt    *zerolog.Event
}

func (t *errEventImpl) Kind() string {
	return "event"
}

func (t *errEventImpl) Event() *zerolog.Event {
	return t.evt
}

func (t *errEventImpl) AddEvent(evt *zerolog.Event) {
	if evt == nil {
		return
	}

	t.evt = evt
}

func (t *errEventImpl) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = t.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, t.String())
	}
}

func (t *errEventImpl) String() string {
	if t.err == nil || generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))

	if t.evt != nil && len(convertEvent(t.evt).buf) > 1 {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorEvent, append(convertEvent(t.evt).buf, '}')))
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, t.caller.String()))
	}

	errStringify(buf, t.err)

	return buf.String()
}

func (t *errEventImpl) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["event"] = json.RawMessage(append(convertEvent(t.evt).buf, '}'))
	return jjson.Marshal(data)
}

func (t *errEventImpl) getData() map[string]any {
	var data = make(map[string]any)
	if t.caller != nil {
		data["caller"] = t.caller
	}

	var mm = errJsonify(t.err)
	if mm != nil {
		for k, v := range mm {
			data[k] = v
		}
	}

	return data
}

func (t *errEventImpl) Unwrap() error {
	return t.err
}

func (t *errEventImpl) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
