package convert

func Map[K comparable, S, D any](src map[K]S, convert func(S) D) map[K]D {
	if src == nil {
		return nil
	}
	dst := make(map[K]D, len(src))
	for k, s := range src {
		dst[k] = convert(s)
	}
	return dst
}

func MapL[A comparable, B any](src []A, convert func(A) B) []B {
	if src == nil {
		return nil
	}

	dst := make([]B, len(src))
	for i := range src {
		dst[i] = convert(src[i])
	}
	return dst
}
