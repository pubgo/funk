package main

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
)

// 单个pkg的error处理

var err1 = &errors.Err{Msg: "业务错误处理", Detail: "详细信息"}

func Hello() {
	defer recovery.Raise(func(err errors.XErr) errors.XErr {
		return err.Wrap("Hello wrap")
	})

	var err2 = errors.WrapF(err1, "处理 wrap")
	assert.MustF(err2, "处理 panic")
	return
}

func CallHello() (gErr error) {
	defer recovery.Recovery(func(err errors.XErr) {
		gErr = err.WrapF("CallHello wrap")
	})

	Hello()

	return
}

func main() {
	defer recovery.Exit()

	assert.Must(CallHello())
}
