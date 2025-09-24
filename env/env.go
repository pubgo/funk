package env

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log/logfields"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/v2/result"
)

func Set(key, value string) result.Error {
	return result.ErrOf(os.Setenv(keyHandler(key), value)).Log(func(e *zerolog.Event) {
		e.Str("key", key)
		e.Str("value", value)
		e.Str(logfields.Msg, "env_set_error")
	})
}

func MustSet(key, value string) { Set(key, value).Must() }

func GetDefault(name string, defaultVal string) string {
	val := Get(name)
	return lo.If(val != "", val).Else(defaultVal)
}

func Get(names ...string) string {
	var val string
	GetVal(&val, names...)
	return val
}

func MustGet(names ...string) string {
	var val string
	GetVal(&val, names...)
	assert.If(val == "", "env not found, names=%q", names)
	return val
}

func GetVal(val *string, names ...string) {
	for _, name := range names {
		env, ok := Lookup(name)
		env = trim(env)
		if ok && env != "" {
			*val = env
			break
		}
	}
}

func GetBoolVal(val *bool, names ...string) {
	dt := Get(names...)
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
	dt := Get(names...)
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
	dt := Get(names...)
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

func Lookup(key string) (string, bool) { return os.LookupEnv(keyHandler(key)) }

func Delete(key string) result.Error {
	return result.ErrOf(os.Unsetenv(keyHandler(key))).Log(func(e *zerolog.Event) {
		e.Str("key", key)
		e.Str(logfields.Msg, "env_delete_error")
	})
}

func Expand(value string) result.Result[string] {
	return result.Wrap(envsubst.String(value)).Log(func(e *zerolog.Event) {
		e.Str("value", value)
		e.Str(logfields.Msg, "env_expand_error")
	})
}

func Map() map[string]string {
	data := make(map[string]string, len(os.Environ()))
	for _, env := range os.Environ() {
		envs := strings.SplitN(env, "=", 2)
		if len(envs) != 2 {
			continue
		}

		data[keyHandler(envs[0])] = envs[1]
	}
	return data
}

func Key(key string) string {
	return keyHandler(key)
}

func LoadFiles(files ...string) (r result.Error) {
	files = lo.Filter(files, func(item string, index int) bool { return pathutil.IsExist(item) })
	if result.Catch(&r, godotenv.Overload(files...)) {
		return
	}

	loadEnv()
	return
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return keyHandler(key), true
}
