package wsclient

import "context"

type WsHandler = func(ctx context.Context, payload Message) Message

type Client interface {
	Start()
	Send(ctx context.Context, payload Message) error
	RegisterAction(act string, handler WsHandler) error
	Close()
}

type Message interface {
	GetID() string
	Marshal() []byte
	GetAction() string
	GetData() interface{}
	ReplyMessage(v interface{}) Message
}
