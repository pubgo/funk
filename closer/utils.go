package closer

import (
	"io"
	"log/slog"
)

func SafeClose(closer io.Closer) {
	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil {
		slog.Error("failed to close operation", "err", err)
	}
}
