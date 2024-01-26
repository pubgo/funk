package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/result"
)

var trim = strings.TrimSpace

func Set(key, value string) error {
	k, v, ok := Normalize(fmt.Sprintf("%s=%s", key, value))
	assert.If(!ok, "env key is incorrect")
	return os.Setenv(Key(k), v)
}

func Get(names ...string) string {
	var val string
	GetWith(&val, names...)
	return trim(val)
}

func MustGet(names ...string) string {
	var val string
	GetWith(&val, names...)
	assert.If(val == "", "env not found, names=%q", names)
	return trim(val)
}

func GetWith(val *string, names ...string) {
	for _, name := range names {
		env, ok := Lookup(name)
		env = trim(env)
		if ok && env != "" {
			*val = trim(env)
			break
		}
	}
}

func GetBoolVal(val *bool, names ...string) {
	var dt = trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.ParseBool(dt)
	if err != nil {
		log.Printf("env: failed to parse string to bool, err=%v\n", err)
		return
	}

	*val = v
}

func GetIntVal(val *int, names ...string) {
	var dt = trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.Atoi(dt)
	if err != nil {
		log.Printf("env: failed to parse string to int, err=%v\n", err)
		return
	}

	*val = v
}

func GetFloatVal(val *float64, names ...string) {
	var dt = trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.ParseFloat(dt, 32)
	if err != nil {
		log.Printf("env: failed to parse string to float, err=%v\n", err)
		return
	}

	*val = v
}

func Lookup(key string) (string, bool) {
	return os.LookupEnv(Key(key))
}

func Delete(key string) error {
	return os.Unsetenv(Key(key))
}

func Expand(value string) result.Result[string] {

	return result.Of(envsubst.String(value))
}

func Map() map[string]string {
	var data = make(map[string]string, len(os.Environ()))
	for _, env := range os.Environ() {
		envs := strings.SplitN(env, "=", 2)
		data[envs[0]] = envs[1]
	}
	return data
}

func Key(key string) string {
	return strings.ToUpper(trim(key))
}
