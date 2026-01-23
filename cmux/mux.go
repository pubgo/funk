// Package cmux provides a compatibility wrapper around connmux.
//
// This repository historically carried experimental code under a "cmux" folder.
// New development should use package connmux instead.
package cmux

import (
	"net"
	"time"

	"github.com/pubgo/funk/v2/connmux"
)

type Matcher = connmux.Matcher

type MatchWriter = connmux.MatchWriter

type Option = connmux.Option

type Mux = connmux.Mux

var ErrServerClosed = connmux.ErrServerClosed

var ErrNotMatched = connmux.ErrNotMatched

func New(l net.Listener, opts ...Option) *Mux { return connmux.New(l, opts...) }

func WithReadTimeout(d time.Duration) Option { return connmux.WithReadTimeout(d) }

func WithMaxSniffBytes(n int) Option { return connmux.WithMaxSniffBytes(n) }

func WithConnBacklog(n int) Option { return connmux.WithConnBacklog(n) }

func WithErrorHandler(h func(error) bool) Option { return connmux.WithErrorHandler(h) }

func Any() Matcher { return connmux.Any() }

func Prefix(prefixes ...[]byte) Matcher { return connmux.Prefix(prefixes...) }

func HTTP1Fast() Matcher { return connmux.HTTP1Fast() }

func HTTP1() Matcher { return connmux.HTTP1() }

func HTTP2() Matcher { return connmux.HTTP2() }

func HTTP2HeaderField(name, value string) Matcher { return connmux.HTTP2HeaderField(name, value) }

func HTTP2HeaderFieldPrefix(name, valuePrefix string) Matcher {
	return connmux.HTTP2HeaderFieldPrefix(name, valuePrefix)
}

func HTTP2HeaderFieldSendSettings(name, value string) MatchWriter {
	return connmux.HTTP2HeaderFieldSendSettings(name, value)
}
