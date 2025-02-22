package helpers

import "iter"

func SliceMap[T any, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))
	for i, v := range s {
		res[i] = f(v)
	}

	return res
}

func ReverseIter[T any](s []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}
