package logger

import (
	"context"
	"github.com/rs/xid"
)

var logCtxKey = xid.New().String()

func NewWithLog() {

}

func GetFromCtx(ctx context.Context) {

}
