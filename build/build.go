package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog/log"
)

func Run(BuildFlags, Packages string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", "-c", "go build "+BuildFlags+" "+Packages)
	cmd.Dir = pwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	log.Printf("go build cmd is: %v", cmd.Args)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}

	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}

	log.Info().Msgf("Go build exit successful.")
	return nil
}

type Options struct {
	Name    string
	Path    string
	Ext     string
	Target  string
	Goos    string
	Goarch  string
	Goamd64 string
	Goarm   string
	Gomips  string
}

// Build contains the build configuration section.
type Build struct {
	BuildDetails          `yaml:",inline" json:",inline"`
	ID                    string                 `yaml:"id,omitempty" json:"id,omitempty"`
	Goos                  []string               `yaml:"goos,omitempty" json:"goos,omitempty"`
	Goarch                []string               `yaml:"goarch,omitempty" json:"goarch,omitempty"`
	Goarm                 []string               `yaml:"goarm,omitempty" json:"goarm,omitempty"`
	Gomips                []string               `yaml:"gomips,omitempty" json:"gomips,omitempty"`
	Goamd64               []string               `yaml:"goamd64,omitempty" json:"goamd64,omitempty"`
	Targets               []string               `yaml:"targets,omitempty" json:"targets,omitempty"`
	Ignore                []IgnoredBuild         `yaml:"ignore,omitempty" json:"ignore,omitempty"`
	Dir                   string                 `yaml:"dir,omitempty" json:"dir,omitempty"`
	Main                  string                 `yaml:"main,omitempty" json:"main,omitempty"`
	Binary                string                 `yaml:"binary,omitempty" json:"binary,omitempty"`
	Hooks                 BuildHookConfig        `yaml:"hooks,omitempty" json:"hooks,omitempty"`
	Builder               string                 `yaml:"builder,omitempty" json:"builder,omitempty"`
	ModTimestamp          string                 `yaml:"mod_timestamp,omitempty" json:"mod_timestamp,omitempty"`
	Skip                  bool                   `yaml:"skip,omitempty" json:"skip,omitempty" jsonschema:"oneof_type=string;boolean"`
	GoBinary              string                 `yaml:"gobinary,omitempty" json:"gobinary,omitempty"`
	Command               string                 `yaml:"command,omitempty" json:"command,omitempty"`
	NoUniqueDistDir       bool                   `yaml:"no_unique_dist_dir,omitempty" json:"no_unique_dist_dir,omitempty"`
	NoMainCheck           bool                   `yaml:"no_main_check,omitempty" json:"no_main_check,omitempty"`
	UnproxiedMain         string                 `yaml:"-" json:"-"` // used by gomod.proxy
	UnproxiedDir          string                 `yaml:"-" json:"-"` // used by gomod.proxy
	BuildDetailsOverrides []BuildDetailsOverride `yaml:"overrides,omitempty" json:"overrides,omitempty"`
}

// IgnoredBuild represents a build ignored by the user.
type IgnoredBuild struct {
	Goos    string `yaml:"goos,omitempty" json:"goos,omitempty"`
	Goarch  string `yaml:"goarch,omitempty" json:"goarch,omitempty"`
	Goarm   string `yaml:"goarm,omitempty" json:"goarm,omitempty" jsonschema:"oneof_type=string;integer"`
	Gomips  string `yaml:"gomips,omitempty" json:"gomips,omitempty"`
	Goamd64 string `yaml:"goamd64,omitempty" json:"goamd64,omitempty"`
}

type BuildDetailsOverride struct {
	Goos         string `yaml:"goos,omitempty" json:"goos,omitempty"`
	Goarch       string `yaml:"goarch,omitempty" json:"goarch,omitempty"`
	Goarm        string `yaml:"goarm,omitempty" json:"goarm,omitempty" jsonschema:"oneof_type=string;integer"`
	Gomips       string `yaml:"gomips,omitempty" json:"gomips,omitempty"`
	Goamd64      string `yaml:"goamd64,omitempty" json:"goamd64,omitempty"`
	BuildDetails `yaml:",inline" json:",inline"`
}

type BuildDetails struct {
	Buildmode string      `yaml:"buildmode,omitempty" json:"buildmode,omitempty"`
	Ldflags   StringArray `yaml:"ldflags,omitempty" json:"ldflags,omitempty"`
	Tags      FlagArray   `yaml:"tags,omitempty" json:"tags,omitempty"`
	Flags     FlagArray   `yaml:"flags,omitempty" json:"flags,omitempty"`
	Asmflags  StringArray `yaml:"asmflags,omitempty" json:"asmflags,omitempty"`
	Gcflags   StringArray `yaml:"gcflags,omitempty" json:"gcflags,omitempty"`
	Env       []string    `yaml:"env,omitempty" json:"env,omitempty"`
}

// StringArray is a wrapper for an array of strings.
type StringArray []string

// UnmarshalYAML is a custom unmarshaler that wraps strings in arrays.
func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strings []string
	if err := unmarshal(&strings); err != nil {
		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}
		*a = []string{str}
	} else {
		*a = strings
	}
	return nil
}

func (a StringArray) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{{
			Type: "string",
		}, {
			Type: "array",
			Items: &jsonschema.Schema{
				Type: "string",
			},
		}},
	}
}

// FlagArray is a wrapper for an array of strings.
type FlagArray []string

// UnmarshalYAML is a custom unmarshaler that wraps strings in arrays.
func (a *FlagArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var flags []string
	if err := unmarshal(&flags); err != nil {
		var flagstr string
		if err := unmarshal(&flagstr); err != nil {
			return err
		}
		*a = strings.Fields(flagstr)
	} else {
		*a = flags
	}
	return nil
}

func (a FlagArray) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{{
			Type: "string",
		}, {
			Type: "array",
			Items: &jsonschema.Schema{
				Type: "string",
			},
		}},
	}
}

type BuildHookConfig struct {
	Pre  Hooks `yaml:"pre,omitempty" json:"pre,omitempty"`
	Post Hooks `yaml:"post,omitempty" json:"post,omitempty"`
}

type Hooks []Hook

// UnmarshalYAML is a custom unmarshaler that allows simplified declaration of single command.
func (bhc *Hooks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var singleCmd string
	err := unmarshal(&singleCmd)
	if err == nil {
		*bhc = []Hook{{Cmd: singleCmd}}
		return nil
	}

	type t Hooks
	var hooks t
	if err := unmarshal(&hooks); err != nil {
		return err
	}
	*bhc = (Hooks)(hooks)
	return nil
}

type Hook struct {
	Dir    string   `yaml:"dir,omitempty" json:"dir,omitempty"`
	Cmd    string   `yaml:"cmd,omitempty" json:"cmd,omitempty"`
	Env    []string `yaml:"env,omitempty" json:"env,omitempty"`
	Output bool     `yaml:"output,omitempty" json:"output,omitempty"`
}

// UnmarshalYAML is a custom unmarshaler that allows simplified declarations of commands as strings.
func (bh *Hook) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cmd string
	if err := unmarshal(&cmd); err != nil {
		type t Hook
		var hook t
		if err := unmarshal(&hook); err != nil {
			return err
		}
		*bh = (Hook)(hook)
		return nil
	}

	bh.Cmd = cmd
	return nil
}

func (bh Hook) JSONSchema() *jsonschema.Schema {
	type t Hook
	reflector := jsonschema.Reflector{
		ExpandedStruct: true,
	}
	schema := reflector.Reflect(&t{})
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type: "string",
			},
			schema,
		},
	}
}
