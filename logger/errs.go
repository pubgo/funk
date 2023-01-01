package logger

import (
	"errors"
	"fmt"
)

var ErrUnknownLevel = errors.New("slog: unknown level name")
var ErrTypeNotMatch = fmt.Errorf("slog: data type not match")
