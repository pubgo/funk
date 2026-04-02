package config

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"dario.cat/mergo"
	"github.com/a8m/envsubst"
	"github.com/samber/lo"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/pathutil"
	"github.com/pubgo/funk/v2/result"
)

// findConfigPath searches for config file and returns path and directory.
// Returns error if config file is not found.
func findConfigPath(name, typ string, configDirs ...string) (cfgPath string, cfgDir string, err error) {
	if len(configDirs) == 0 {
		configDirs = append(configDirs, "./", defaultConfigPath)
	}

	if name == "" {
		name = defaultConfigName
	}

	if typ == "" {
		typ = defaultConfigType
	}

	configName := fmt.Sprintf("%s.%s", name, typ)
	var notFoundPath []string
	for _, path := range getPathList() {
		for _, dir := range configDirs {
			cfgPath := filepath.Join(path, dir, configName)
			if pathutil.IsNotExist(cfgPath) {
				notFoundPath = append(notFoundPath, cfgPath)
			} else {
				return cfgPath, filepath.Dir(cfgPath), nil
			}
		}
	}

	return "", "", errors.Errorf("config not found in: %v", notFoundPath)
}

// getPathList 递归得到当前目录到跟目录中所有的目录路径
//
//	paths: [./, ../, ../../, ..., /]
func getPathList() (paths []string) {
	wd := assert.Must1(filepath.Abs(""))
	for len(wd) > 0 && !os.IsPathSeparator(wd[len(wd)-1]) {
		paths = append(paths, wd)
		wd = filepath.Dir(wd)
	}
	return paths
}

func MergeR[A any, B any | *any](dst *A, src ...B) (r result.Result[*A]) {
	if len(src) == 0 {
		return r.WithValue(dst)
	}

	err := Merge(dst, src...)
	if err != nil {
		return r.WithErr(err)
	}
	return r.WithValue(dst)
}

func Merge[A any, B any | *any](dst *A, src ...B) error {
	for i := range src {
		err := mergo.Merge(
			dst,
			src[i],
			mergo.WithOverride,
			mergo.WithAppendSlice,
			mergo.WithTransformers(new(transformer)),
		)
		if err != nil {
			return errors.WrapTags(err, errors.Tags{
				"dst_type": reflect.TypeOf(dst).String(),
				"src_type": reflect.TypeOf(src[i]).String(),
				"dst":      dst,
				"src":      src[i],
			})
		}
	}
	return nil
}

type transformer struct{}

func (s *transformer) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if t == nil || t.Kind() != reflect.Slice {
		return nil
	}

	if !t.Elem().Implements(reflect.TypeOf((*NamedConfig)(nil)).Elem()) {
		return nil
	}

	return func(dst, src reflect.Value) error {
		if !src.IsValid() || src.IsNil() {
			return nil
		}

		dstMap := make(map[string]any)
		for i := 0; i < dst.Len(); i++ {
			c := dst.Index(i).Interface()
			dstMap[c.(NamedConfig).ConfigUniqueName()] = c
		}

		for i := 0; i < src.Len(); i++ {
			c := src.Index(i).Interface()
			uniqueName := c.(NamedConfig).ConfigUniqueName()
			if dstMap[uniqueName] == nil {
				dstMap[uniqueName] = c
				continue
			}

			d := dstMap[uniqueName]
			err := mergo.Merge(d, c, mergo.WithOverride, mergo.WithAppendSlice, mergo.WithTransformers(new(transformer)))
			if err != nil {
				return errors.WrapFn(err, func() errors.Tags {
					return errors.Tags{
						"dst":      d,
						"src":      c,
						"src-type": reflect.TypeOf(c).String(),
						"dst-type": reflect.TypeOf(d).String(),
					}
				})
			}
		}

		data := lo.MapToSlice(dstMap, func(key string, value any) reflect.Value { return reflect.ValueOf(value) })
		dst.Set(makeList(dst.Type().Elem(), data))
		return nil
	}
}

func unmarshalOneOrList[T any](list *[]T, value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		var t T
		if err := value.Decode(&t); err != nil {
			return err
		}
		*list = append(*list, t)
		return nil
	}

	if value.Kind == yaml.SequenceNode {
		return value.Decode(list)
	}
	return errors.Errorf("unmarshalled node: %v", value.Value)
}

func listAllPath(dirOrPath string) (ret result.Result[[]string]) {
	if !pathutil.IsDir(dirOrPath) {
		return ret.WithValue([]string{dirOrPath})
	}

	var paths []string
	walk := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		paths = append(paths, path)
		return nil
	}
	err := filepath.Walk(dirOrPath, walk)
	if err != nil {
		return ret.WithErr(err)
	}
	return ret.WithValue(paths)
}

func makeList(typ reflect.Type, data []reflect.Value) reflect.Value {
	val := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(data))
	return reflect.Append(val, data...)
}

type config struct {
	workDir    string
	envSpecMap EnvSpecMap // allowed env vars from patch_envs
}

// RegisterExpr registers a custom expression function for use in config templates.
// Returns error if the name already exists. For backward compatibility, use MustRegisterExpr for panic behavior.
// Note: Custom functions must have simple signatures: func() T or func(T) R
func RegisterExpr(name string, fn any) error {
	return globalManager.RegisterExprFunc(name, fn)
}

func evalData(template []byte, cfg *config) []byte {
	cleanedTemplate := removeYAMLComments(template)

	engine, engineErr := newCelEngine(cfg)

	exprTpl := fasttemplate.New(string(cleanedTemplate), "${{", "}}")
	res := []byte(exprTpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		tag = strings.TrimSpace(tag)
		if engineErr != nil {
			return -1, errors.Wrap(engineErr, "failed to create CEL engine")
		}

		d, err := result.WrapErr(engine.Eval(tag))
		if err.IsErr() {
			err.Log(func(e result.Event) {
				e.Str("tag", tag)
			})
			return -1, err.Err()
		}

		data, err := result.WrapErr(yaml.Marshal(d))
		if err.IsErr() {
			err.Log(func(e result.Event) {
				e.Str("tag", tag)
				e.Msg("failed to marshal yaml")
			})
			return -1, err.Err()
		}

		return w.Write(bytes.TrimSpace(data))
	}))

	envTpl := fasttemplate.New(string(res), "${", "}")
	return []byte(envTpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		tag = strings.TrimSpace(tag)
		name := strings.ToUpper(strings.TrimSpace(strings.Split(tag, ":")[0]))
		if cfg.envSpecMap != nil {
			if _, defined := cfg.envSpecMap[name]; !defined {
				return -1, fmt.Errorf("env: variable %q is not defined in envs, all env vars must be declared", name)
			}
		}

		tag = fmt.Sprintf("${%s}", tag)
		return w.Write(result.Wrap(envsubst.Bytes([]byte(tag))).
			Map(bytes.TrimSpace).
			UnwrapOrLog(func(e result.Event) {
				e.Str("env", name)
				e.Msg("failed to process env subst")
			}))
	}))
}

// removeYAMLCommentsFromLine removes comments from a YAML line while respecting quoted strings
func removeYAMLCommentsFromLine(line []byte) []byte {
	resultData := make([]byte, 0, len(line))
	inSingleQuote := false
	inDoubleQuote := false
	i := 0
	for i < len(line) {
		char := line[i]

		// Check for escape character (backslash)
		if char == '\\' && (inSingleQuote || inDoubleQuote) {
			// In YAML, within single quotes, backslash has no special meaning
			// Within double quotes, backslash can escape certain characters
			if inDoubleQuote && i+1 < len(line) {
				// Check if next character is a quote or backslash
				nextChar := line[i+1]
				if nextChar == '"' || nextChar == '\\' {
					// This is an escaped quote or backslash, keep both characters
					resultData = append(resultData, char, nextChar)
					i += 2
					continue
				}
			}
			// For single quotes or other cases, just append the backslash
			resultData = append(resultData, char)
			i++
			continue
		}

		// Check for quote characters, but not if escaped (handled above)
		if char == '\'' && !inDoubleQuote {
			// Toggle single quote state
			inSingleQuote = !inSingleQuote
			resultData = append(resultData, char)
		} else if char == '"' && !inSingleQuote {
			// Toggle double quote state
			inDoubleQuote = !inDoubleQuote
			resultData = append(resultData, char)
		} else if char == '#' && !inSingleQuote && !inDoubleQuote {
			// Found comment marker outside of quotes, stop processing
			break
		} else {
			resultData = append(resultData, char)
		}
		i++
	}
	// Trim trailing spaces
	return bytes.TrimRight(resultData, " \t")
}

// removeYAMLComments removes all comments from YAML data while respecting quoted strings
func removeYAMLComments(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	var cleanedLines [][]byte
	for _, line := range lines {
		// Check if original line is empty (only whitespace)
		originalTrimmed := bytes.TrimSpace(line)
		isOriginalEmpty := len(originalTrimmed) == 0

		// Process each line to remove comments
		cleanedLine := removeYAMLCommentsFromLine(line)
		trimmed := bytes.TrimSpace(cleanedLine)

		// Preserve empty lines, but remove lines that were only comments
		if isOriginalEmpty {
			// Original line was empty, preserve it
			cleanedLines = append(cleanedLines, cleanedLine)
		} else if len(trimmed) > 0 {
			// Line had content and still has content after comment removal
			cleanedLines = append(cleanedLines, cleanedLine)
		}
		// If original line had content but after comment removal it's empty,
		// it means the line was only a comment, so we skip it
	}
	return bytes.Join(cleanedLines, []byte("\n"))
}

// evalExpr evaluates a CEL expression with the given config context
func evalExpr(code string, cfg *config) (any, error) {
	engine, err := newCelEngine(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CEL engine")
	}
	return engine.Eval(code)
}
