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
	"github.com/samber/lo"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
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

// MustRegisterExpr is like RegisterExpr but panics on error.
// Deprecated: prefer RegisterExpr which returns error.
func MustRegisterExpr(name string, fn any) {
	if err := RegisterExpr(name, fn); err != nil {
		panic(err)
	}
}

func cfgFormat(template []byte, cfg *config) []byte {
	tpl := fasttemplate.New(string(template), "${{", "}}")
	return []byte(tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		tag = strings.TrimSpace(tag)
		evalData, err := eval(tag, cfg)
		if err != nil {
			return -1, errors.Wrap(err, tag)
		}

		data, err := yaml.Marshal(evalData)
		if err != nil {
			log.Err(err).
				Str("tag", tag).
				Msgf("failed to marshal yaml: %v", evalData)
			return -1, errors.Wrap(err, tag)
		}

		return w.Write(bytes.TrimSpace(data))
	}))
}

// eval evaluates a CEL expression with the given config context
func eval(code string, cfg *config) (any, error) {
	engine, err := newCelEngine(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CEL engine")
	}
	return engine.Eval(code)
}
