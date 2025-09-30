package shutil

import (
	"bytes"
	"fmt"
	result2 "github.com/pubgo/funk/v2/result"
	"os"
	"os/exec"
	"strings"

	"github.com/pubgo/funk/v2/log/logfields"
	"github.com/rs/zerolog"
)

func Run(args ...string) (r result2.Result[string]) {
	defer result2.Recovery(&r)

	b := bytes.NewBufferString("")

	cmd := Shell(args...)
	cmd.Stdout = b

	result2.ErrOf(cmd.Run()).Must(func(e *zerolog.Event) {
		e.Str(logfields.Msg, fmt.Sprintf("failed to execute: "+strings.Join(args, " ")))
	})

	return r.WithValue(strings.TrimSpace(b.String()))
}

func GoModGraph() result2.Result[string] {
	return Run("go", "mod", "graph")
}

func GoList() result2.Result[string] {
	return Run("go", "list", "./...")
}

func GraphViz(in, out string) (err error) {
	ret := Run("dot", "-Tsvg", in)
	if ret.IsErr() {
		return ret.GetErr()
	}

	return os.WriteFile(out, []byte(ret.GetValue()), 0o600)
}

func Shell(args ...string) *exec.Cmd {
	shell := strings.Join(args, " ")
	cmd := exec.Command("/bin/sh", "-c", shell)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd
}
