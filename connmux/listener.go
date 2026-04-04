package connmux

import (
	"net"
	"sync"

	"github.com/pubgo/funk/v2/errors"
)

var errListenerClosed = errors.New("connmux: listener closed")

// muxListener is a net.Listener backed by an internal connection queue.
// It receives connections that were matched by the parent Mux.
//
// It intentionally does NOT close accepted conns when it is closed; it only
// stops future Accept calls.
type muxListener struct {
	addr  net.Addr
	connc chan net.Conn

	donec chan struct{}
	once  sync.Once
}

func newMuxListener(addr net.Addr, backlog int) *muxListener {
	if backlog <= 0 {
		backlog = 128
	}
	return &muxListener{
		addr:  addr,
		connc: make(chan net.Conn, backlog),
		donec: make(chan struct{}),
	}
}

func (l *muxListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.connc:
		if c == nil {
			return nil, errors.Wrap(errListenerClosed, "nil connection received")
		}
		return c, nil
	case <-l.donec:
		return nil, errors.Wrap(errListenerClosed, "listener closed during accept")
	}
}

func (l *muxListener) Close() error {
	l.once.Do(func() {
		close(l.donec)
		close(l.connc)
	})
	return nil
}

func (l *muxListener) Addr() net.Addr { return l.addr }

func (l *muxListener) enqueue(c net.Conn, muxDone <-chan struct{}) error {
	select {
	case <-l.donec:
		return errors.Wrap(errListenerClosed, "listener closed")
	case <-muxDone:
		return errors.Wrap(ErrServerClosed, "server closed during enqueue")
	case l.connc <- c:
		return nil
	}
}
