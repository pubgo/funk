package config

import (
	"bytes"
	_ "embed"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/a8m/envsubst"
	"github.com/pubgo/funk/env"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type testCfg struct {
	Assets struct {
		TestMd struct {
			TestAbc struct {
				Secret Base64File `yaml:"secret"`
			} `yaml:"test_abc"`
		} `yaml:"test_md"`
	} `yaml:"assets"`
}

//go:embed configs/assets/.gen.yaml
var genYaml string

func TestExpr(t *testing.T) {
	os.Setenv("testAbc", "hello")
	env.Init()

	assert.Equal(t, string(cfgFormat([]byte("{{env.TEST_ABC}}"), &config{})), "hello")
	assert.Equal(t, string(cfgFormat([]byte(`{{embed("configs/assets/secret")}}`), &config{})), strings.TrimSpace(`MTIzNDU2CjEyMzQ1NgoxMjM0NTYKMTIzNDU2CjEyMzQ1NgoxMjM0NTYKMTIzNDU2CjEyMzQ1Ng==`))

	var dd, err = os.ReadFile("configs/assets/assets.yaml")
	assert.NoError(t, err)
	var dd1 = bytes.TrimSpace(cfgFormat(dd, &config{workDir: "configs/assets"}))
	var cfg testCfg
	assert.NoError(t, yaml.Unmarshal(dd1, &cfg))

	assert.Equal(t, string(dd1), strings.TrimSpace(genYaml))
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
	sort.Slice(cfg.Names, func(i, j int) bool {
		return cfg.Names[i].Name < cfg.Names[j].Name
	})

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
