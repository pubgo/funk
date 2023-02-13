package lifecycle

func New(handlers []Handler) Lifecycle {
	var lc = new(lifecycleImpl)
	for i := range handlers {
		handlers[i](lc)
	}
	return lc
}
