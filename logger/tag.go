package logger

import _ "golang.org/x/exp/slog"

type tag struct {
	key   string
	value interface{}
}

func (t *tag) Key() string {
	return t.key
}

func (t *tag) Value() interface{} {
	return t.value
}

func Tag(key string, val interface{}) Tagger {
	return &tag{key: key, value: val}
}

func Tags(key string, val ...interface{}) Tagger {
	return &tag{key: key, value: val}
}
