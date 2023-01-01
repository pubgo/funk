package logger

import (
	"errors"
	"fmt"
	"strings"
)

const (
	UNKNOWN Level = iota
	TRACE
	DEBUG
	INFO
	WARNING
	ERROR
	CRITICAL
)

// Level holds a severity level.
type Level uint8

// ParseLevel converts a string representation of a logging level to a
// Level. It returns the level and whether it was valid or not.
func ParseLevel(level string) (Level, error) {
	level = strings.ToUpper(level)
	switch level {
	case "TRACE":
		return TRACE, nil
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN", "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "CRITICAL":
		return CRITICAL, nil
	default:
		return UNKNOWN, fmt.Errorf("failed to parse log level, err=%w", ErrUnknownLevel)
	}
}

func (level *Level) Enabled(ll Level) bool {
	return ll <= *level
}

// String implements Stringer.
func (level *Level) String() string {
	switch *level {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Short returns a five character string to use in
// aligned logging output.
func (level *Level) Short() string {
	switch *level {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case WARNING:
		return "WARN "
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITC"
	default:
		return "     "
	}
}

func (level *Level) UnmarshalJSON(text []byte) error {
	if level == nil {
		return errors.New("can't unmarshal a nil *Level")
	}

	var err error
	*level, err = ParseLevel(string(text))
	if err != nil {
		return fmt.Errorf("can't unmarshal %s, err=%w", text, err)
	}

	return nil
}

func (level *Level) MarshalJSON() ([]byte, error) {
	return []byte(level.String()), nil
}
