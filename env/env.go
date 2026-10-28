package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/pathutil"
	"github.com/pubgo/funk/v2/result"
)

const Name = "env"

func Set(key, value string) result.Error {
	return result.ErrOf(os.Setenv(keyHandler(key), value)).
		Log(func(e result.Event) {
			e.Str("key", key)
			e.Str("value", value)
			e.Msg("env_set_error")
		})
}

func MustSet(key, value string) { Set(key, value).MustWithLog() }

func Get(names ...string) string {
	var val string
	getVal(&val, names...)
	return val
}

func MustGet(names ...string) string {
	val := Get(names...)
	assert.If(val == "", "env value not set, names=%q", names)
	return val
}

func GetOr(name, defaultVal string) string {
	val := Get(name)
	return lo.If(val != "", val).Else(defaultVal)
}

func GetWith(val *string, names ...string) { getVal(val, names...) }

func getVal(val *string, names ...string) {
	for _, name := range names {
		v, ok := lookupEnv(name)
		if !ok || trim(v) == "" {
			continue
		}

		*val = v
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
		getLog().Error("failed to parse string to bool", "keys", names, "value", val, "err", err)
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
		getLog().Error("failed to parse string to int", "keys", names, "value", val, "err", err)
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
		getLog().Error("failed to parse string to float", "keys", names, "value", val, "err", err)
		return -1
	}

	return v
}

// GetBoolOr returns the parsed bool value of name, defaultVal when the
// variable is unset, empty or not a valid bool.
func GetBoolOr(name string, defaultVal bool) bool {
	var val string
	getVal(&val, name)
	if val == "" {
		return defaultVal
	}

	v, err := strconv.ParseBool(val)
	if err != nil {
		getLog().Error("failed to parse string to bool", "key", name, "value", val, "err", err)
		return defaultVal
	}

	return v
}

// GetIntOr returns the parsed int value of name, defaultVal when the variable
// is unset, empty or not a valid int.
func GetIntOr(name string, defaultVal int) int {
	var val string
	getVal(&val, name)
	if val == "" {
		return defaultVal
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		getLog().Error("failed to parse string to int", "key", name, "value", val, "err", err)
		return defaultVal
	}

	return v
}

// GetFloatOr returns the parsed float value of name, defaultVal when the
// variable is unset, empty or not a valid float.
func GetFloatOr(name string, defaultVal float64) float64 {
	var val string
	getVal(&val, name)
	if val == "" {
		return defaultVal
	}

	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		getLog().Error("failed to parse string to float", "key", name, "value", val, "err", err)
		return defaultVal
	}

	return v
}

func Lookup(key string) (string, bool) { return lookupEnv(key) }

// Delete unsets every stored spelling of key (raw, normalized and normalized
// uppercase), so variables saved under a renamed key are removed as well.
func Delete(key string) result.Error {
	return result.ErrOf(unsetKey(key)).
		Log(func(e result.Event) {
			e.Str("key", key)
			e.Msg("env_delete_error")
		})
}

func unsetKey(key string) error {
	var lastErr error
	seen := make(map[string]struct{}, 3)
	for _, k := range envKeyCandidates(key) {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}

		if err := os.Unsetenv(k); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func MustDelete(key string) { Delete(key).MustWithLog() }

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
		return r
	}

	var needReloadEnv bool
	for _, file := range files {
		data, err := result.WrapErr(os.ReadFile(file))
		if err.Throw(&r) {
			r.Log(func(e result.Event) { e.Msgf("failed to read file %q", file) })
			return
		}

		dataMap, err := result.WrapErr(godotenv.UnmarshalBytes(data))
		if err.Throw(&r) {
			r.Log(func(e result.Event) { e.Msgf("failed to parse env file:%s", file) })
			return
		}

		for k, v := range dataMap {
			if k == "" || v == "" {
				continue
			}

			if Set(k, v).Throw(&r) {
				return r
			}

			needReloadEnv = true
		}
	}

	if needReloadEnv {
		loadEnv()
	}

	return r
}

// Normalize a-b=>a_b, a.b=>a_b, a/b=>a_b
func Normalize(key string) (string, bool) {
	key = trim(key)
	if key == "" || strings.HasPrefix(key, "_") || strings.HasPrefix(key, "=") {
		return key, false
	}

	return keyHandler(key), true
}
