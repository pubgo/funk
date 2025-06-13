package shutil

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/v2/result"
)

func Run(args ...string) (r result.Result[string]) {
	defer result.RecoveryErr(&r)

	b := bytes.NewBufferString("")

	cmd := Shell(args...)
	cmd.Stdout = b

	err := result.ErrOf(cmd.Run()).
		Inspect(func(err error) {
			log.Err(err).Msg("failed to execute: " + strings.Join(args, " "))
		}).
		CatchErr(&r)
	if err {
		return
	}

	return r.WithValue(strings.TrimSpace(b.String()))
}

func GoModGraph() result.Result[string] {
	return Run("go", "mod", "graph")
}

func GoList() result.Result[string] {
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
