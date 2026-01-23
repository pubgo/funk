package connmux_test

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pubgo/funk/v2/closer"
	"github.com/pubgo/funk/v2/connmux"
)

func ExampleMux() {
	root, _ := net.Listen("tcp", "127.0.0.1:0")
	defer closer.SafeClose(root)

	m := connmux.New(root, connmux.WithReadTimeout(2*time.Second))
	httpL := m.Match(connmux.HTTP1Fast())
	_ = m.Match(connmux.Any())

	go func() { _ = m.Serve() }()
	defer closer.SafeClose(m)

	c, _ := net.Dial("tcp", root.Addr().String())
	defer closer.SafeClose(c)
	_, _ = c.Write([]byte("GET / HTTP/1.1\r\nHost: example\r\n\r\n"))

	s, _ := httpL.Accept()
	defer closer.SafeClose(s)

	b := make([]byte, 3)
	_, _ = io.ReadFull(s, b)
	fmt.Println(string(b))

	// Output:
	// GET
}
