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

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/log/logfields"
	"github.com/pubgo/funk/v2/pathutil"
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

func Get(names ...string) string {
	var val string
	getVal(&val, names...)
	return val
}

func MustGet(names ...string) string {
	val := Get(names...)
	assert.If(val == "", "env not found, names=%q", names)
	return val
}

func GetOr(name string, defaultVal string) string {
	val := Get(name)
	return lo.If(val != "", val).Else(defaultVal)
}

func GetWith(val *string, names ...string) { getVal(val, names...) }

func getVal(val *string, names ...string) {
	for _, name := range names {
		env := trim(os.Getenv(keyHandler(name)))
		if env == "" {
			continue
		}

		*val = env
		break
	}
}

func GetBool(names ...string) bool {
	var val string
	getVal(&val, names...)
	if val == "" {
		return false
	}

	v, err := strconv.ParseBool(val)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to bool, keys=%q value=%s err=%v", names, val, err))
		return false
	}

	return v
}

func GetInt(names ...string) int {
	var val string
	getVal(&val, names...)
	if val == "" {
		return -1
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to int, keys=%q value=%s err=%v", names, val, err))
		return -1
	}

	return v
}

func GetFloat(names ...string) float64 {
	var val string
	getVal(&val, names...)
	if val == "" {
		return -1
	}

	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("env: failed to parse string to float, keys=%q value=%s err=%v", names, val, err))
		return -1
	}

	return v
}

func Lookup(key string) (string, bool) { return os.LookupEnv(keyHandler(key)) }

func Delete(key string) result.Error {
	return result.ErrOf(os.Unsetenv(keyHandler(key))).Log(func(e *zerolog.Event) {
		e.Str("key", key)
		e.Str(logfields.Msg, "env_delete_error")
	})
}

func MustDelete(key string) { Delete(key).Must() }

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
	if len(files) == 0 {
		return
	}

	for _, file := range files {
		data := result.Wrap(os.ReadFile(file)).Unwrap(&r)
		if r.IsErr() {
			return
		}

		dataMap := result.Wrap(godotenv.UnmarshalBytes(data)).Unwrap(&r)
		if r.IsErr() {
			return
		}

		for k, v := range dataMap {
			if k == "" || v == "" {
				continue
			}

			if Set(k, v).Catch(&r) {
				return
			}
		}
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
