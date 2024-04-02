package typex

type Ctx map[string]any

func (c Ctx) ToTuple() Tuple {
	tt := make(Tuple, 0, len(c))
	for k := range c {
		tt = append(tt, KV{K: k, V: c[k]})
	}
	return tt
}

type (
	List  []any
	Tuple []KV
)

func (t Tuple) ToCtx() Ctx {
	ctx := make(Ctx, len(t))
	for i := range t {
		ctx[t[i].K] = t[i].V
	}
	return ctx
}

type KV struct {
	K string `json:"key"`
	V any    `json:"value"`
}

func StrOf(s ...string) []string {
	return s
}
