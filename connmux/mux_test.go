package connmux

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"github.com/pubgo/funk/v2/closer"
)

func startMux(t *testing.T, m *Mux) {
	t.Helper()
	go func() {
		_ = m.Serve()
	}()
	t.Cleanup(func() { _ = m.Close() })
}

func dialAndWrite(t *testing.T, addr string, b []byte) {
	t.Helper()
	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer closer.SafeClose(c)
	_ = c.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if _, err := c.Write(b); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func acceptOne(t *testing.T, l net.Listener) net.Conn {
	t.Helper()
	type deadlineListener interface{ SetDeadline(time.Time) error }
	if dl, ok := l.(deadlineListener); ok {
		_ = dl.SetDeadline(time.Now().Add(3 * time.Second))
	}

	c, err := l.Accept()
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	return c
}

func TestMux_PrefixAndAnyDispatch(t *testing.T) {
	root, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	m := New(root, WithReadTimeout(2*time.Second))
	lA := m.Match(Prefix([]byte("A")))
	lAny := m.Match(Any())
	startMux(t, m)

	dialAndWrite(t, root.Addr().String(), []byte("Ahello"))
	c1 := acceptOne(t, lA)
	defer closer.SafeClose(c1)
	buf := make([]byte, 6)
	_ = c1.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.ReadFull(c1, buf); err != nil {
		t.Fatalf("read1: %v", err)
	}
	if string(buf) != "Ahello" {
		t.Fatalf("got %q", string(buf))
	}

	dialAndWrite(t, root.Addr().String(), []byte("Bhello"))
	c2 := acceptOne(t, lAny)
	defer closer.SafeClose(c2)
	buf2 := make([]byte, 6)
	_ = c2.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.ReadFull(c2, buf2); err != nil {
		t.Fatalf("read2: %v", err)
	}
	if string(buf2) != "Bhello" {
		t.Fatalf("got %q", string(buf2))
	}
}

func TestMux_HTTP2HeaderFieldDispatch_ReplaysPreface(t *testing.T) {
	root, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	m := New(root, WithReadTimeout(2*time.Second))
	grpcL := m.Match(HTTP2HeaderField("content-type", "application/grpc"))
	_ = m.Match(HTTP2())
	startMux(t, m)

	c, err := net.Dial("tcp", root.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer closer.SafeClose(c)
	_ = c.SetDeadline(time.Now().Add(3 * time.Second))

	// Send client preface.
	if _, err := c.Write([]byte(http2.ClientPreface)); err != nil {
		t.Fatalf("write preface: %v", err)
	}

	// Send a single HEADERS frame with content-type.
	var hb bytes.Buffer
	enc := hpack.NewEncoder(&hb)
	_ = enc.WriteField(hpack.HeaderField{Name: ":method", Value: "POST"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":scheme", Value: "http"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":path", Value: "/grpc"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":authority", Value: "example"})
	_ = enc.WriteField(hpack.HeaderField{Name: "content-type", Value: "application/grpc"})

	fr := http2.NewFramer(c, c)
	if err := fr.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		BlockFragment: hb.Bytes(),
		EndStream:     false,
		EndHeaders:    true,
	}); err != nil {
		t.Fatalf("write headers: %v", err)
	}

	s := acceptOne(t, grpcL)
	defer closer.SafeClose(s)
	pref := make([]byte, len(http2.ClientPreface))
	_ = s.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.ReadFull(s, pref); err != nil {
		t.Fatalf("server read preface: %v", err)
	}
	if string(pref) != http2.ClientPreface {
		t.Fatalf("preface mismatch")
	}
}

func TestMux_HTTP2HeaderFieldSendSettings(t *testing.T) {
	root, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	m := New(root, WithReadTimeout(2*time.Second))
	grpcL := m.MatchWithWriters(HTTP2HeaderFieldSendSettings("content-type", "application/grpc"))
	_ = m.Match(Any())
	startMux(t, m)

	c, err := net.Dial("tcp", root.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer closer.SafeClose(c)
	_ = c.SetDeadline(time.Now().Add(4 * time.Second))

	if _, err := c.Write([]byte(http2.ClientPreface)); err != nil {
		t.Fatalf("write preface: %v", err)
	}

	fr := http2.NewFramer(c, c)
	if err := fr.WriteSettings(); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	// Expect the mux to respond with SETTINGS before we send HEADERS.
	f, err := fr.ReadFrame()
	if err != nil {
		t.Fatalf("read frame: %v", err)
	}
	if _, ok := f.(*http2.SettingsFrame); !ok {
		t.Fatalf("expected settings frame, got %T", f)
	}

	var hb bytes.Buffer
	enc := hpack.NewEncoder(&hb)
	_ = enc.WriteField(hpack.HeaderField{Name: ":method", Value: "POST"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":scheme", Value: "http"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":path", Value: "/grpc"})
	_ = enc.WriteField(hpack.HeaderField{Name: ":authority", Value: "example"})
	_ = enc.WriteField(hpack.HeaderField{Name: "content-type", Value: "application/grpc"})

	if err := fr.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		BlockFragment: hb.Bytes(),
		EndStream:     false,
		EndHeaders:    true,
	}); err != nil {
		t.Fatalf("write headers: %v", err)
	}

	s := acceptOne(t, grpcL)
	defer closer.SafeClose(s)
	// Ensure stream is intact: preface should be readable.
	pref := make([]byte, len(http2.ClientPreface))
	_ = s.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.ReadFull(s, pref); err != nil {
		t.Fatalf("server read preface: %v", err)
	}
}

func TestMux_CloseUnblocksAccept(t *testing.T) {
	root, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	m := New(root)
	l := m.Match(Any())
	startMux(t, m)

	if err := m.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	_, err = l.Accept()
	if err == nil {
		t.Fatalf("expected accept error after close")
	}
}

func TestMux_SniffOverflowIsFatal(t *testing.T) {
	root, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	var (
		gotErrMu sync.Mutex
		gotErr   error
	)
	m := New(root,
		WithMaxSniffBytes(8),
		WithReadTimeout(2*time.Second),
		WithErrorHandler(func(err error) bool {
			gotErrMu.Lock()
			gotErr = err
			gotErrMu.Unlock()
			return true
		}),
	)
	// First matcher will try to read 9 bytes and trigger overflow.
	_ = m.Match(Prefix([]byte("123456789")))
	anyL := m.Match(Any())
	startMux(t, m)

	dialAndWrite(t, root.Addr().String(), []byte("123456789"))

	// Any should NOT receive this connection because overflow is fatal.
	// Accept should block until we close the mux.
	acc := make(chan error, 1)
	go func() {
		c, err := anyL.Accept()
		if c != nil {
			_ = c.Close()
		}
		acc <- err
	}()

	select {
	case err := <-acc:
		// If it returns without us closing the mux, that's a failure.
		t.Fatalf("unexpected accept return: %v", err)
	case <-time.After(200 * time.Millisecond):
		// expected
	}

	_ = m.Close()
	select {
	case err := <-acc:
		if err == nil {
			t.Fatalf("expected accept error after close")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("accept did not unblock after close")
	}
	gotErrMu.Lock()
	handled := gotErr
	gotErrMu.Unlock()
	if handled == nil {
		t.Fatalf("expected error handler to be called")
	}
}
