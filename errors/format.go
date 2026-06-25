package errors

import (
	"fmt"
	"strings"
)

func layerMessage(err error) string {
	switch e := err.(type) {
	case *ErrWrap:
		if msg, ok := e.Tags["msg"]; ok {
			return strings.TrimSpace(fmt.Sprint(msg))
		}
	case *Err:
		return strings.TrimSpace(e.Msg)
	case Err:
		return strings.TrimSpace(e.Msg)
	}

	if err == nil {
		return ""
	}
	return strings.TrimSpace(err.Error())
}

// FormatChain returns a human-readable error chain joined with ": ".
// Wrap context from ErrWrap tags is included while duplicate adjacent messages are omitted.
func FormatChain(err error) string {
	if err == nil {
		return ""
	}

	parts := make([]string, 0, 4)
	Walk(err, func(current error) bool {
		msg := layerMessage(current)
		if msg == "" {
			return true
		}

		if len(parts) == 0 || parts[len(parts)-1] != msg {
			parts = append(parts, msg)
		}
		return true
	})

	return strings.Join(parts, ": ")
}

// FullMessage is an alias for FormatChain.
func FullMessage(err error) string {
	return FormatChain(err)
}
