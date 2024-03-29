package cmux

//import (
//	"fmt"
//	"net"
//	"net/http"
//	"os"
//	"strings"
//
//	"github.com/pubgo/funk/config"
//	"github.com/pubgo/funk/log"
//	"github.com/soheilhy/cmux"
//	"github.com/tmc/grpc-websocket-proxy/wsproxy"
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"go.etcd.io/etcd/client/v3/naming/resolver"
//	"golang.org/x/net/http2"
//	"google.golang.org/grpc"
//	// https://github.com/shaxbee/go-wsproxy
//	"go.etcd.io/etcd/client/pkg/v3/transport"
//	_ "go.etcd.io/etcd/client/v3/naming/resolver"
//	"go.etcd.io/etcd/pkg/v3/httputil"
//)
//
//func init() {
//	cli, cerr := clientv3.NewFromURL("http://localhost:2379")
//	etcdResolver, err := resolver.NewBuilder(cli)
//	conn, gerr := grpc.Dial("etcd:///foo/bar/my-service", grpc.WithResolvers(etcdResolver))
//
//	wsproxy.WebsocketProxy(
//		gwmux,
//		wsproxy.WithRequestMutator(
//			// Default to the POST method for streams
//			func(_ *http.Request, outgoing *http.Request) *http.Request {
//				outgoing.Method = "POST"
//				return outgoing
//			},
//		),
//	)
//
//	host := httputil.GetHostname(req)
//
//	m := cmux.New(sctx.l)
//	grpcl := m.Match(cmux.HTTP2())
//	go func() { errHandler(gs.Serve(grpcl)) }()
//
//	httpl := m.Match(cmux.HTTP1())
//	go func() { errHandler(srvhttp.Serve(httpl)) }()
//
//	var tlsl net.Listener
//	tlsl, err = transport.NewTLSListener(m.Match(cmux.Any()), tlsinfo)
//	if err != nil {
//		return err
//	}
//
//	m.Serve()
//}
//
//func configureHttpServer(srv *http.Server, cfg config.ServerConfig) error {
//	// todo (ahrtr): should we support configuring other parameters in the future as well?
//	return http2.ConfigureServer(srv, &http2.Server{
//		MaxConcurrentStreams: cfg.MaxConcurrentStreams,
//	})
//}
//
//// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
//// connections or otherHandler otherwise. Given in gRPC docs.
//func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
//	if otherHandler == nil {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			grpcServer.ServeHTTP(w, r)
//		})
//	}
//
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
//			grpcServer.ServeHTTP(w, r)
//		} else {
//			otherHandler.ServeHTTP(w, r)
//		}
//	})
//}
//
//// GetHostname returns the hostname from request Host field.
//// It returns empty string, if Host field contains invalid
//// value (e.g. "localhost:::" with too many colons).
//func GetHostname(req *http.Request) string {
//	if req == nil {
//		return ""
//	}
//
//	h, _, err := net.SplitHostPort(req.Host)
//	if err != nil {
//		return req.Host
//	}
//	return h
//}
//
//func mustListenCMux(lg log.Logger, tlsinfo *transport.TLSInfo) cmux.CMux {
//	l, err := net.Listen("tcp", grpcProxyListenAddr)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//
//	if l, err = transport.NewKeepAliveListener(l, "tcp", nil); err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//	if tlsinfo != nil {
//		tlsinfo.CRLFile = grpcProxyListenCRL
//		if l, err = transport.NewTLSListener(l, tlsinfo); err != nil {
//			lg.Err(err).Msg("failed to create TLS listener")
//		}
//	}
//
//	lg.Info().Str("address", grpcProxyListenAddr).Msg("listening for gRPC proxy client requests")
//	return cmux.New(l)
//}
