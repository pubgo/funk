package connmux

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/pubgo/funk/v2/closer"
	"github.com/pubgo/funk/v2/errors"
)

var errSniffOverflow = errors.New("connmux: sniff buffer overflow")

// muxConn wraps a net.Conn and records bytes read during sniffing, so that
// they can be replayed to the selected protocol server.
type muxConn struct {
	net.Conn

	mu sync.Mutex

	buf []byte
	// servePos is used once sniffing is stopped; it replays buf to the server.
	servePos int

	maxSniffBytes int
	readTimeout   time.Duration

	sniffStopped bool
	sniffErr     error
}

func newMuxConn(c net.Conn, maxSniffBytes int, readTimeout time.Duration) *muxConn {
	if maxSniffBytes <= 0 {
		maxSniffBytes = 1 << 20
	}
	return &muxConn{
		Conn:          c,
		maxSniffBytes: maxSniffBytes,
		readTimeout:   readTimeout,
	}
}

func (c *muxConn) stopSniffing() {
	c.mu.Lock()
	c.sniffStopped = true
	c.servePos = 0
	c.mu.Unlock()
}

func (c *muxConn) sniffError() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return errors.Wrap(c.sniffErr, "sniff error")
}

func (c *muxConn) Read(p []byte) (int, error) {
	c.mu.Lock()
	if c.sniffStopped {
		// Replay buffered bytes first.
		if c.servePos < len(c.buf) {
			n := copy(p, c.buf[c.servePos:])
			c.servePos += n
			c.mu.Unlock()
			return n, nil
		}
		c.mu.Unlock()
		n, err := c.Conn.Read(p)
		return n, errors.Wrap(err, "read from conn")
	}
	c.mu.Unlock()

	// During sniffing, the Mux uses sniffReader() which does its own buffering.
	// If someone reads directly from the conn before dispatch, we still buffer.
	n, err := c.readFromSource(p)
	return n, errors.Wrap(err, "read from source")
}

func (c *muxConn) readFromSource(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.readTimeout > 0 {
		err := c.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			return 0, errors.Wrap(err, "set read deadline")
		}

		defer closer.ErrClose(func() error {
			return c.SetReadDeadline(time.Time{})
		})
	}

	n, err := c.Conn.Read(p)
	if n > 0 {
		if len(c.buf)+n > c.maxSniffBytes {
			// We cannot drop bytes (would corrupt stream), so fail fast.
			c.sniffErr = errSniffOverflow
			return 0, errors.Wrap(errSniffOverflow, "sniff buffer overflow")
		}
		c.buf = append(c.buf, p[:n]...)
	}
	return n, errors.Wrap(err, "read from underlying conn")
}

// sniffReader returns a fresh reader starting from the beginning of the
// sniffed stream. Reads beyond the current buffer are pulled from the
// underlying connection and appended to the buffer.
func (c *muxConn) sniffReader() io.Reader {
	return &sniffReader{c: c}
}

type sniffReader struct {
	c   *muxConn
	pos int
}

func (r *sniffReader) Read(p []byte) (int, error) {
	r.c.mu.Lock()
	// Replay already-buffered bytes first.
	if r.pos < len(r.c.buf) {
		n := copy(p, r.c.buf[r.pos:])
		r.pos += n
		r.c.mu.Unlock()
		return n, nil
	}
	r.c.mu.Unlock()

	// Need to read more from the source.
	n, err := r.c.readFromSource(p)
	if n > 0 {
		r.pos += n
	}
	return n, errors.Wrap(err, "sniff read from source")
}

var _ net.Conn = (*muxConn)(nil)
var _ io.Reader = (*sniffReader)(nil)
