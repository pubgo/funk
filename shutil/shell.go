package shutil

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pubgo/funk/errors/errcheck"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/recovery"
	"github.com/pubgo/funk/result"
)

func Run(args ...string) (r result.Result[string]) {
	defer recovery.Err(&r.E)

	b := bytes.NewBufferString("")

	cmd := Shell(args...)
	cmd.Stdout = b

	err := cmd.Run()
	errcheck.Inspect(err, func(err error) {
		log.Err(err).Msg("fail to execute: " + strings.Join(args, " "))
	})
	if errcheck.Check(&r.E, err) {
		return
	}

	return r.WithVal(strings.TrimSpace(b.String()))
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
