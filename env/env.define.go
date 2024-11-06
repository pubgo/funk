package env

import (
	"strings"

	"github.com/pubgo/funk/stack"
	"github.com/samber/lo"
)

var envs []*DefinedEnv

func Define(key string, desc ...string) *DefinedEnv {
	var e = &DefinedEnv{key: KeyHandler(key), desc: strings.Join(desc, " "), stack: stack.Caller(1)}
	envs = append(envs, e)
	return e
}

func FindAllDefinedEnvs() []*DefinedEnvWrap {
	return lo.Map(envs, func(env *DefinedEnv, index int) *DefinedEnvWrap {
		return &DefinedEnvWrap{
			Key:   env.key,
			Desc:  env.desc,
			Stack: env.stack,
		}
	})
}

type DefinedEnvWrap struct {
	Key   string
	Desc  string
	Stack *stack.Frame
}

type DefinedEnv struct {
	key   string
	desc  string
	stack *stack.Frame
}

func (e DefinedEnv) Get() string     { return Get(e.key) }
func (e DefinedEnv) MustGet() string { return MustGet(e.key) }
func (e DefinedEnv) Bool() (v bool) {
	GetBoolVal(&v, e.key)
	return
}
