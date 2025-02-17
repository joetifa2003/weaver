package helpers

func SliceMap[T any, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))
	for i, v := range s {
		res[i] = f(v)
	}

	return res
}
