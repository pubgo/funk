package connmux

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/pubgo/funk/v2/errors"
)

// Matcher matches a connection by reading from r.
//
// Notes:
//   - r is a sniffing reader: reads are buffered and will be replayed to the
//     final protocol server once the connection is dispatched.
//   - Each matcher is evaluated against a fresh reader starting from the
//     beginning of the buffered stream (OR semantics).
type Matcher func(r io.Reader) bool

// MatchWriter is like Matcher but can also write back to the connection
// during sniffing (e.g. sending HTTP/2 SETTINGS to unblock certain clients).
type MatchWriter func(w io.Writer, r io.Reader) bool

// ErrServerClosed is returned by (*Mux).Serve after Close is called.
var ErrServerClosed = errors.New("connmux: server closed")

// ErrNotMatched indicates that an accepted connection did not match any rule.
var ErrNotMatched = errors.New("connmux: connection not matched")

// Option configures a Mux.
type Option func(*Mux)

// WithReadTimeout sets a per-read deadline while sniffing.
// A zero duration disables deadlines (default).
func WithReadTimeout(d time.Duration) Option {
	return func(m *Mux) { m.readTimeout = d }
}

// WithMaxSniffBytes caps how many bytes can be buffered while matching.
// Defaults to 1 MiB.
func WithMaxSniffBytes(n int) Option {
	return func(m *Mux) {
		if n > 0 {
			m.maxSniffBytes = n
		}
	}
}

// WithConnBacklog sets the per-matched-listener connection backlog.
// Defaults to 128.
func WithConnBacklog(n int) Option {
	return func(m *Mux) {
		if n > 0 {
			m.backlog = n
		}
	}
}

// WithErrorHandler sets a handler for accept/match errors.
// If h returns true, Serve continues; otherwise Serve returns the error.
func WithErrorHandler(h func(error) bool) Option {
	return func(m *Mux) {
		if h != nil {
			m.errh = h
		}
	}
}

// Mux multiplexes a single net.Listener into multiple protocol-specific listeners.
type Mux struct {
	root net.Listener

	mu    sync.RWMutex
	rules []*rule

	donec chan struct{}
	once  sync.Once

	readTimeout   time.Duration
	maxSniffBytes int
	backlog       int
	errh          func(error) bool
}

type rule struct {
	l        *muxListener
	matchers []Matcher
	writers  []MatchWriter
}

// New creates a new connection multiplexer.
func New(l net.Listener, opts ...Option) *Mux {
	m := &Mux{
		root:          l,
		donec:         make(chan struct{}),
		readTimeout:   0,
		maxSniffBytes: 1 << 20,
		backlog:       128,
		errh:          func(error) bool { return true },
	}
	for _, opt := range opts {
		if opt != nil {
			opt(m)
		}
	}
	return m
}

// Match registers a rule and returns a listener that only accepts matching connections.
// Match order defines priority.
func (m *Mux) Match(matchers ...Matcher) net.Listener {
	ml := newMuxListener(m.root.Addr(), m.backlog)
	m.mu.Lock()
	m.rules = append(m.rules, &rule{l: ml, matchers: append([]Matcher(nil), matchers...)})
	m.mu.Unlock()
	return ml
}

// MatchWithWriters registers a rule composed of match-writers.
// Match order defines priority.
func (m *Mux) MatchWithWriters(writers ...MatchWriter) net.Listener {
	ml := newMuxListener(m.root.Addr(), m.backlog)
	m.mu.Lock()
	m.rules = append(m.rules, &rule{l: ml, writers: append([]MatchWriter(nil), writers...)})
	m.mu.Unlock()
	return ml
}

// Serve starts accepting on the root listener and dispatching to matched listeners.
// It blocks until the root listener errors or Close is called.
func (m *Mux) Serve() error {
	for {
		c, err := m.root.Accept()
		if err != nil {
			select {
			case <-m.donec:
				return errors.Wrap(ErrServerClosed, "server closed")
			default:
			}
			werr := errors.Wrap(err, "accept error")
			if m.errh != nil && m.errh(werr) {
				continue
			}
			return werr
		}
		go m.dispatch(c)
	}
}

// Close stops Serve and closes the root listener and all matched listeners.
func (m *Mux) Close() error {
	m.once.Do(func() {
		close(m.donec)
		_ = m.root.Close()
		m.mu.RLock()
		rules := append([]*rule(nil), m.rules...)
		m.mu.RUnlock()
		for _, r := range rules {
			_ = r.l.Close()
		}
	})
	return nil
}

func (m *Mux) dispatch(c net.Conn) {
	mc := newMuxConn(c, m.maxSniffBytes, m.readTimeout)

	m.mu.RLock()
	rules := append([]*rule(nil), m.rules...)
	m.mu.RUnlock()

	for _, r := range rules {
		if r == nil || r.l == nil {
			continue
		}

		matched := false
		for _, mm := range r.matchers {
			if mm == nil {
				continue
			}
			r := mc.sniffReader()
			if mm(r) {
				matched = true
				break
			}
			if serr := mc.sniffError(); serr != nil {
				if m.errh != nil {
					_ = m.errh(errors.Wrap(serr, "connmux: sniff"))
				}
				_ = mc.Close()
				return
			}
		}
		if !matched {
			for _, mw := range r.writers {
				if mw == nil {
					continue
				}
				r := mc.sniffReader()
				if mw(mc, r) {
					matched = true
					break
				}
				if serr := mc.sniffError(); serr != nil {
					if m.errh != nil {
						_ = m.errh(errors.Wrap(serr, "connmux: sniff"))
					}
					_ = mc.Close()
					return
				}
			}
		}

		if !matched {
			continue
		}

		mc.stopSniffing()
		if err := r.l.enqueue(mc, m.donec); err != nil {
			if m.errh != nil {
				_ = m.errh(errors.Wrap(err, "connmux: enqueue"))
			}
			_ = mc.Close()
		}
		return
	}

	if m.errh != nil {
		_ = m.errh(ErrNotMatched)
	}
	_ = mc.Close()
}
