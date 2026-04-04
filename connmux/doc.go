// Package connmux multiplexes network connections based on their initial bytes.
//
// It lets you serve multiple protocols (gRPC/HTTP/1/HTTP/2/raw TCP, etc.) on the
// same listening port by sniffing the beginning of each accepted connection and
// dispatching it to a protocol-specific net.Listener.
//
// This package is conceptually similar to github.com/soheilhy/cmux (Apache-2.0)
// but is an independent implementation tailored for this repository.
//
// ## Basics
//
//	root, _ := net.Listen("tcp", ":8080")
//	m := connmux.New(root,
//		connmux.WithReadTimeout(2*time.Second),
//		connmux.WithMaxSniffBytes(1<<20),
//	)
//
//	// Match order defines priority.
//	grpcL := m.Match(connmux.HTTP2HeaderField("content-type", "application/grpc"))
//	httpL := m.Match(connmux.HTTP1Fast())
//	other := m.Match(connmux.Any())
//
//	go grpcServer.Serve(grpcL)
//	go httpServer.Serve(httpL)
//	go serveRaw(other)
//	_ = m.Serve()
//
// ## MatchWithWriters (Java gRPC)
//
// Some clients (notably Java gRPC) may wait for the server SETTINGS frame before
// sending request headers. In those cases use MatchWithWriters with
// HTTP2HeaderFieldSendSettings:
//
//	grpcL := m.MatchWithWriters(
//		connmux.HTTP2HeaderFieldSendSettings("content-type", "application/grpc"),
//	)
//
// ## Notes / limitations
//
//   - The match decision is made when a connection is accepted. A single
//     connection cannot switch protocols later.
//   - Matching is based on reading and buffering initial bytes. Use
//     WithMaxSniffBytes to cap memory per connection and WithReadTimeout to
//     avoid slowloris-style hangs during sniffing.
//   - If you terminate TLS after connmux, some stdlib components may not detect
//     the underlying *tls.Conn due to type assertions on net.Conn wrappers.
//     If your handler relies on TLS-specific state, prefer terminating TLS
//     before multiplexing.
package connmux
