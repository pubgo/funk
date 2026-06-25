package main

import (
	"fmt"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errcode"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/proto/errorpb"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func main() {
	defer recovery.Exit()

	demoBasicErrors()
	demoResultAndLog()
	demoErrCode()
}

func demoBasicErrors() {
	err := errors.New("database connection failed", errors.Tags{
		"component": "database",
		"host":      "localhost:5432",
	})
	err = errors.Wrap(err, "initialize payment service")

	fmt.Println("root message:", err.Error())
	fmt.Println("full chain:", errors.FormatChain(err))
	fmt.Println("user tags:", errors.CollectUserTags(err))
	fmt.Println("error id:", errors.GetErrorId(err))
}

func demoResultAndLog() {
	r := result.ErrOf(errors.Wrap(
		errors.New("not found", errors.Tags{"user_id": 42}),
		"load user",
	))

	if r.IsErr() {
		log.Error().
			Str("message", r.Message()).
			Interface("tags", r.Tags()).
			Err(r.Err()).
			Msg("operation failed")
	}
}

func demoErrCode() {
	const name = "example.demo.not_found"
	_ = errcode.RegisterErrCode(&errorpb.ErrCode{
		Code:       404,
		Message:    "resource not found",
		Name:       name,
		StatusCode: errorpb.Code_NotFound,
	})

	code, ok := errcode.LookupErrCode(name)
	if !ok {
		panic("missing registered error code")
	}

	err := errcode.NewCodeErr(code)
	fmt.Println(errors.FormatChain(err))
}
