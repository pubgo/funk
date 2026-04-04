package connmux

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

// HTTP1Fast matches common HTTP/1.x methods using a small prefix check.
// It is faster but less accurate than HTTP1().
func HTTP1Fast() Matcher {
	// Include "PRI " to avoid confusing HTTP/2 preface with HTTP/1.
	prefixes := [][]byte{
		[]byte("GET "),
		[]byte("POST "),
		[]byte("PUT "),
		[]byte("PATCH "),
		[]byte("DELETE "),
		[]byte("HEAD "),
		[]byte("OPTIONS "),
		[]byte("CONNECT "),
		[]byte("TRACE "),
		[]byte("PRI "),
	}
	return Prefix(prefixes...)
}

// HTTP1 matches by parsing an HTTP/1.x request line and headers.
func HTTP1() Matcher {
	return func(r io.Reader) bool {
		br := bufio.NewReader(r)
		req, err := http.ReadRequest(br)
		if err != nil {
			return false
		}
		_ = req.Body.Close()
		return true
	}
}

// HTTP2 matches an HTTP/2 connection by verifying the client preface.
func HTTP2() Matcher {
	return func(r io.Reader) bool {
		return hasHTTP2Preface(r)
	}
}

// HTTP2HeaderField matches an HTTP/2 connection by looking for an exact header
// field value (case-insensitive name match).
//
// Note: some clients (notably Java gRPC) will not send HEADERS until they
// receive a SETTINGS frame from the server; for those, use
// HTTP2HeaderFieldSendSettings with MatchWithWriters.
func HTTP2HeaderField(name, value string) Matcher {
	name = strings.ToLower(name)
	return func(r io.Reader) bool {
		return matchHTTP2HeaderField(nil, r, name, func(v string) bool { return v == value })
	}
}

// HTTP2HeaderFieldPrefix matches an HTTP/2 connection by looking for a header
// field value prefix (case-insensitive name match).
func HTTP2HeaderFieldPrefix(name, valuePrefix string) Matcher {
	name = strings.ToLower(name)
	return func(r io.Reader) bool {
		return matchHTTP2HeaderField(nil, r, name, func(v string) bool { return strings.HasPrefix(v, valuePrefix) })
	}
}

// HTTP2HeaderFieldSendSettings is a MatchWriter that writes an initial SETTINGS
// frame (when appropriate) while sniffing, then matches based on a header field.
//
// This is useful for clients that wait for server SETTINGS before sending
// request HEADERS (e.g. Java gRPC).
func HTTP2HeaderFieldSendSettings(name, value string) MatchWriter {
	name = strings.ToLower(name)
	return func(w io.Writer, r io.Reader) bool {
		return matchHTTP2HeaderField(w, r, name, func(v string) bool { return v == value })
	}
}

func hasHTTP2Preface(r io.Reader) bool {
	preface := []byte(http2.ClientPreface)
	buf := make([]byte, len(preface))
	if _, err := io.ReadFull(r, buf); err != nil {
		return false
	}
	return bytes.Equal(buf, preface)
}

func matchHTTP2HeaderField(w io.Writer, r io.Reader, name string, matches func(string) bool) bool {
	if !hasHTTP2Preface(r) {
		return false
	}

	var out io.Writer
	if w != nil {
		out = w
	} else {
		out = io.Discard
	}

	fr := http2.NewFramer(out, r)
	matched := false
	foundName := false

	dec := hpack.NewDecoder(4<<10, func(hf hpack.HeaderField) {
		if strings.ToLower(hf.Name) == name {
			foundName = true
			if matches(hf.Value) {
				matched = true
			}
		}
	})

	// Decode until we see the header name (or the stream headers end).
	for {
		f, err := fr.ReadFrame()
		if err != nil {
			return false
		}

		switch f := f.(type) {
		case *http2.SettingsFrame:
			// If the client sent SETTINGS (non-ACK) and we are allowed to write,
			// respond with our own SETTINGS to unblock certain clients.
			if w != nil && !f.IsAck() {
				if err := fr.WriteSettings(); err != nil {
					return false
				}
			}

		case *http2.HeadersFrame:
			if _, err := dec.Write(f.HeaderBlockFragment()); err != nil {
				return false
			}
			if f.HeadersEnded() {
				if foundName {
					return matched
				}
				// Otherwise continue; the header might be on another stream.
			}

		case *http2.ContinuationFrame:
			if _, err := dec.Write(f.HeaderBlockFragment()); err != nil {
				return false
			}
			if f.HeadersEnded() {
				if foundName {
					return matched
				}
			}
		}
	}
}
