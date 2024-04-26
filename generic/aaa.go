package generic

type Ctx[T any] map[string]T

func (c Ctx[T]) ToTuple() Tuple[T] {
	tt := make(Tuple[T], 0, len(c))
	for k := range c {
		tt = append(tt, KV[T]{K: k, V: c[k]})
	}
	return tt
}

type (
	List[T any]  []T
	Tuple[T any] []KV[T]
)

func (t Tuple[T]) ToCtx() Ctx[T] {
	ctx := make(Ctx[T], len(t))
	for i := range t {
		ctx[t[i].K] = t[i].V
	}
	return ctx
}

type KV[T any] struct {
	K string `json:"key"`
	V T      `json:"value"`
}
