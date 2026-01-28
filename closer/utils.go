package closer

import (
	"io"
	"log/slog"
)

func ErrClose(closer func() error) {
	if closer == nil {
		return
	}

	if err := closer(); err != nil {
		slog.Error("failed to close error operation", "err", err)
	}
}

func SafeClose(closer io.Closer) {
	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil {
		slog.Error("failed to close io operation", "err", err)
	}
}
