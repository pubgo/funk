package logger

import "context"

type Tagger interface {
	Key() string
	Value() interface{}
}

type Logger interface {
	Info(level uint, msg string, tags ...Tagger)
	Error(err error, msg string, tags ...Tagger)
}

type Hook interface {
	Hook(tags []Tagger) []Tagger
}

type CtxParser func(ctx context.Context) (bool, []Tagger)
type ValueParser func(v interface{}) (bool, string)
