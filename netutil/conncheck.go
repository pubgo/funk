//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd || solaris || illumos
// +build linux darwin dragonfly freebsd netbsd openbsd solaris illumos

package netutil

import (
	"io"
	"net"
	"syscall"

	"github.com/pubgo/funk/errors"
)

var errUnexpectedRead = errors.New("unexpected read from socket")

func ConnCheck(conn net.Conn) error {
	var sysErr error

	sysConn, ok := conn.(syscall.Conn)
	if !ok {
		return nil
	}

	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err
	}

	err = rawConn.Read(func(fd uintptr) bool {
		var buf [1]byte
		n, err := syscall.Read(int(fd), buf[:])
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = errUnexpectedRead
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}

		return true
	})
	if err != nil {
		return err
	}

	return sysErr
}
