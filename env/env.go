package env

import (
	"fmt"
	"github.com/pubgo/funk/pathutil"
	"github.com/samber/lo"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/v2/result"
)

func Set(key, value string) result.Error {
	return result.ErrOf(os.Setenv(KeyHandler(key), value))
}

func Get(names ...string) string {
	var val string
	GetVal(&val, names...)
	return trim(val)
}

func MustGet(names ...string) string {
	var val string
	GetVal(&val, names...)
	assert.If(val == "", "env not found, names=%q", names)
	return trim(val)
}

func GetVal(val *string, names ...string) {
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
	dt := trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.ParseBool(dt)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to bool, keys=%q value=%s err=%v", names, dt, err))
		return
	}

	*val = v
}

func GetIntVal(val *int, names ...string) {
	dt := trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.Atoi(dt)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to int, keys=%q value=%s err=%v", names, dt, err))
		return
	}

	*val = v
}

func GetFloatVal(val *float64, names ...string) {
	dt := trim(Get(names...))
	if dt == "" {
		return
	}

	v, err := strconv.ParseFloat(dt, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to float, keys=%q value=%s err=%v", names, dt, err))
		return
	}

	*val = v
}

func Lookup(key string) (string, bool) {
	return os.LookupEnv(Key(key))
}

func Delete(key string) result.Error {
	return result.ErrOf(os.Unsetenv(Key(key)))
}

func Expand(value string) result.Result[string] {
	return result.Wrap(envsubst.String(value))
}

func Map() map[string]string {
	data := make(map[string]string, len(os.Environ()))
	for _, env := range os.Environ() {
		envs := strings.SplitN(env, "=", 2)
		data[envs[0]] = envs[1]
	}
	return data
}

func Key(key string) string {
	return KeyHandler(key)
}

func LoadFiles(files ...string) (r result.Error) {
	files = lo.Filter(files, func(item string, index int) bool { return pathutil.IsExist(item) })
	if result.CatchErr(&r, godotenv.Load(files...)) {
		return
	}

	loadEnv()
	return
}
