package config

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"dario.cat/mergo"
	"github.com/expr-lang/expr"
	"github.com/samber/lo"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/result"
)

func GetConfigDir() string {
	return configDir
}

func GetConfigPath() string {
	return configPath
}

func getConfigPath(name, typ string, configDir ...string) (string, string) {
	if len(configDir) == 0 {
		configDir = append(configDir, "./", defaultConfigPath)
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
		for _, dir := range configDir {
			cfgPath := filepath.Join(path, dir, configName)
			if pathutil.IsNotExist(cfgPath) {
				notFoundPath = append(notFoundPath, cfgPath)
			} else {
				return cfgPath, filepath.Dir(cfgPath)
			}
		}
	}

	log.Fatal().Msgf("config not found in: %v", notFoundPath)

	return "", ""
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
	return
}

func SetConfigPath(confPath string) {
	assert.If(confPath == "", "config path is null")
	configPath = confPath
}

func MergeR[A any, B any | *any](dst *A, src ...B) (ret result.Result[*A]) {
	if len(src) == 0 {
		return ret.WithVal(dst)
	}

	err := Merge(dst, src...)
	if err != nil {
		return ret.WithErr(err)
	}
	return ret.WithVal(dst)
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
			return errors.WrapTag(err,
				errors.T("dst_type", reflect.TypeOf(dst).String()),
				errors.T("dst", dst),
				errors.T("src_type", reflect.TypeOf(src[i]).String()),
				errors.T("src", src[i]),
			)
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
			var uniqueName = c.(NamedConfig).ConfigUniqueName()
			if dstMap[uniqueName] == nil {
				dstMap[uniqueName] = c
				continue
			}

			d := dstMap[uniqueName]
			err := mergo.Merge(d, c, mergo.WithOverride, mergo.WithAppendSlice, mergo.WithTransformers(new(transformer)))
			if err != nil {
				return errors.WrapFn(err, func() errors.Tags {
					return errors.Tags{
						errors.T("dst", d),
						errors.T("src", c),
						errors.T("dst-type", reflect.TypeOf(d).String()),
						errors.T("src-type", reflect.TypeOf(c).String()),
					}
				})
			}
		}

		var data = lo.MapToSlice(dstMap, func(key string, value any) reflect.Value { return reflect.ValueOf(value) })
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
	return errors.Format("unmarshalled node: %v", value.Value)
}

func listAllPath(dirOrPath string) (ret result.Result[[]string]) {
	if !pathutil.IsDir(dirOrPath) {
		return ret.WithVal([]string{dirOrPath})
	}

	var paths []string
	var walk = func(path string, info fs.FileInfo, err error) error {
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
	return ret.WithVal(paths)
}

func makeList(typ reflect.Type, data []reflect.Value) reflect.Value {
	val := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	return reflect.Append(val, data...)
}

type config struct {
	workDir string
}

func getEnvData(cfg *config) map[string]any {
	return map[string]any{
		"env": env.Map(),
		"get_path_dir": func() string {
			return cfg.workDir
		},
		"embed": func(name string) string {
			if name == "" {
				return ""
			}

			var path = filepath.Join(cfg.workDir, name)
			var d, err = os.ReadFile(path)
			if err != nil {
				log.Err(err).Str("path", path).Msg("failed to read file")
				return ""
			}

			return strings.TrimSpace(base64.StdEncoding.EncodeToString(d))
		},
	}
}

func cfgFormat(template string, cfg *config) string {
	tpl := fasttemplate.New(template, "${{", "}}")
	return tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		tag = strings.TrimSpace(tag)
		var data, err = yaml.Marshal(eval(tag, cfg))
		if err != nil {
			return -1, errors.Wrap(err, tag)
		}

		return w.Write(bytes.TrimSpace(data))
	})
}

func eval(code string, cfg *config) any {
	envData := getEnvData(cfg)
	data, err := expr.Eval(strings.TrimSpace(code), envData)
	if err != nil {
		panic(err)
	}
	return data
}
