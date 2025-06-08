package shutil

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pubgo/funk/anyhow"
	"github.com/pubgo/funk/assert"
)

func Run(args ...string) (r anyhow.Result[string]) {
	defer anyhow.RecoveryErr(&r.Err)

	b := bytes.NewBufferString("")

	cmd := Shell(args...)
	cmd.Stdout = b

	assert.Must(cmd.Run(), strings.Join(args, " "))
	return r.SetWithValue(strings.TrimSpace(b.String()))
}

func GoModGraph() anyhow.Result[string] {
	return Run("go", "mod", "graph")
}

func GoList() anyhow.Result[string] {
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
