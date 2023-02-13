package logutil

import (
	"io"
	"net/http"
	"strings"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/try"
)

// GracefulClose drains http.Response.Body until it hits EOF
// and closes it. This prevents TCP/TLS connections from closing,
// therefore available for reuse.
// Borrowed from golang/net/context/ctxhttp/cancelreq.go.
func GracefulClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

func HandleClose(log log.Logger, fn func() error) {
	if fn == nil || log == nil {
		log.Error().Msgf("log and fn are all required")
		return
	}

	var err = fn()
	if generic.IsNil(err) {
		return
	}

	log.Err(err).Msg("failed to handle close")
}

func LogOrErr(log log.Logger, msg string, fn func() error) {
	msg = strings.TrimSpace(msg)
	log = log.WithCallerSkip(1)

	var err = try.Try(fn)
	if generic.IsNil(err) {
		log.Info().Msg(msg)
	} else {
		log.Err(err).Msg(msg)
	}
}

func OkOrFailed(log log.Logger, msg string, fn func() error) {
	log = log.WithCallerSkip(1)
	log.Info().Msg(msg)

	var err = try.Try(fn)
	if generic.IsNil(err) {
		log.Info().Msg(msg + " ok")
	} else {
		log.Err(err).Msg(msg + " failed")
	}
}

func ErrRecord(logger log.Logger, err error, fn func(evt *log.Event) string) {
	if generic.IsNil(err) {
		return
	}

	var evt = log.NewEvent()
	var msg = fn(evt)
	logger.WithCallerSkip(1).Err(err).Func(log.WithEvent(evt)).Msg(msg)
}
