package anyhow

func DoBlock(fn func()) { fn() }

func DoBlock1[T any](fn func() T) T { return fn() }
