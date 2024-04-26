package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
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

	log.Panicf("config not found in: %v\n", notFoundPath)

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
	assert.If(configPath == "", "config path is null")
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
		err := mergo.Merge(dst, src[i], mergo.WithOverride, mergo.WithAppendSlice, mergo.WithTransformers(new(transformer)))
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

		dstMap := make(map[string]NamedConfig)
		for i := 0; i < dst.Len(); i++ {
			c := dst.Index(i).Interface().(NamedConfig)
			dstMap[c.ConfigUniqueName()] = c
		}

		for i := 0; i < src.Len(); i++ {
			c := src.Index(i).Interface().(NamedConfig)
			if dstMap[c.ConfigUniqueName()] == nil {
				dst = reflect.Append(dst, reflect.ValueOf(c))
				dstMap[c.ConfigUniqueName()] = c
				continue
			}

			d := dstMap[c.ConfigUniqueName()]
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
