package log

import (
	"github.com/pubgo/funk/logger"
)

type tag struct {
	key   string
	value interface{}
}

func (t *tag) Key() string        { return t.key }
func (t *tag) Value() interface{} { return t.value }

func Tag(key string, val interface{}) logger.Tagger {
	return &tag{key: key, value: val}
}

func Tags(key string, val ...interface{}) logger.Tagger {
	if len(val) == 0 {
		return nil
	}

	return &tag{key: key, value: val}
}
