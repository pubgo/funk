package logger

import (
	"reflect"
)

type tag struct {
	key   string
	value interface{}
}

func (t *tag) Kind() reflect.Kind {

}

func (t *tag) Key() string        { return t.key }
func (t *tag) Value() interface{} { return t.value }

func Tag(key string, val interface{}) Tagger {
	return &tag{key: key, value: val}
}

func Tags(key string, val ...interface{}) Tagger {
	if len(val) == 0 {
		return &tag{key: key, value: "[]"}
	}

	return &tag{key: key, value: val}
}
