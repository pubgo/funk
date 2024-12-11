package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a8m/envsubst"
	expr "github.com/expr-lang/expr"
	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/log"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"
)

func init() {
	//expr.Function()
}

var envData = map[string]interface{}{
	"env": env.Map(),
	"embed": func(name string, dir string) string {
		if name == "" {
			return ""
		}

		var path = filepath.Join(dir, name)
		var d, err = os.ReadFile(path)
		if err != nil {
			log.Err(err).Str("path", path).Msg("failed to read file")
			return ""
		}

		return strings.TrimSpace(string(d))
	},
}

func init() {
	fmt.Printf("%#v\n", envData)
}

func eval(code string, dir string) any {
	data, err := expr.Eval(strings.TrimSpace(code), envData)
	if err != nil {
		panic(err)
	}
	return data
}

func TestExpr(t *testing.T) {
	t.Log(Format("${{env.SHELL}}", ""))
	t.Log(Format(`${{embed("configs/assets/secret","")}}`, ""))

	var dd, err = os.ReadFile("configs/assets/assets.yaml")
	assert.NoError(t, err)
	t.Log("\n" + Format(string(dd), "configs/assets"))
	os.WriteFile("configs/assets/.gen.yaml", []byte(Format(string(dd), "configs/assets")), 0644)
}

func Format(template string, dir string) string {
	tpl := fasttemplate.New(template, "${{", "}}")
	return tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		fmt.Println(tag)
		var data, err = yaml.Marshal(eval(tag, dir))
		if err != nil {
			return -1, err
		}
		return w.Write(data)
	})
}

func TestEnv(t *testing.T) {
	os.Setenv("hello", "world")
	data, err := envsubst.String("${hello}")
	assert.Nil(t, err)
	assert.Equal(t, data, "world")

	os.Setenv("hello", "")
	data, err = envsubst.String("${hello:-abc}")
	assert.Nil(t, err)
	assert.Equal(t, data, "abc")

	data, err = envsubst.String("${{hello:-abc}}")
	assert.Nil(t, err)
	assert.Equal(t, data, "${{hello:-abc}}")
}

func TestConfigPath(t *testing.T) {
	t.Log(getConfigPath("", ""))
	assert.Panics(t, func() {
		t.Log(getConfigPath("", "toml"))
	})
}

var _ NamedConfig = (*configL)(nil)

type configL struct {
	Name  string
	Value string
}

func (c configL) ConfigUniqueName() string {
	return c.Name
}

type configA struct {
	Names []*configL
	Name1 configL
}

func TestMerge(t *testing.T) {
	cfg := &configA{}
	assert.Nil(t, Merge(
		cfg,
		configA{
			Name1: configL{
				Name: "a1",
			},
		},
	))
	assert.Equal(t, cfg.Name1.Name, "a1")

	cfg = &configA{}
	assert.Nil(t, Merge(
		cfg,
		configA{
			Name1: configL{
				Name: "a1",
			},
		},
		configA{
			Names: []*configL{
				{Name: "a2"},
			},
			Name1: configL{
				Name: "a2",
			},
		},
	))
	assert.Equal(t, cfg.Name1.Name, "a2")
	assert.Equal(t, len(cfg.Names), 1)
	assert.Equal(t, cfg.Names[0].Name, "a2")

	cfg = new(configA)
	assert.Nil(t, Merge(
		cfg,
		configA{
			Name1: configL{
				Name: "a1",
			},
		},

		configA{
			Names: []*configL{
				{Name: "a2", Value: "a2"},
			},
			Name1: configL{
				Name: "a2",
			},
		},

		configA{
			Names: []*configL{
				{Name: "a2", Value: "a3"},
				{Name: "a3"},
			},
			Name1: configL{
				Name: "a3",
			},
		},
	))
	assert.Equal(t, cfg.Name1.Name, "a3")
	assert.Equal(t, len(cfg.Names), 2)
	assert.Equal(t, cfg.Names[0].Name, "a2")
	assert.Equal(t, cfg.Names[0].Value, "a3")
	assert.Equal(t, cfg.Names[1].Name, "a3")
	assert.Equal(t, cfg.Names[1].Value, "")

	cfg = new(configA)
	assert.Nil(t, Merge(
		cfg,
		configA{
			Name1: configL{
				Name:  "a1",
				Value: "a1",
			},
		},

		configA{
			Names: []*configL{
				{Name: "a1", Value: ""},
			},
			Name1: configL{
				Name: "a1",
			},
		},
	))
	assert.Equal(t, cfg.Name1.Name, "a1")
	assert.Equal(t, cfg.Name1.Value, "a1")
	assert.Equal(t, len(cfg.Names), 1)
	assert.Equal(t, cfg.Names[0].Name, "a1")
	assert.Equal(t, cfg.Names[0].Value, "")
}
