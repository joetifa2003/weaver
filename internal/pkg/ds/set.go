package ds

import "iter"

type Set[T comparable] struct {
	data map[T]struct{}
}

func (s *Set[T]) Add(v T) {
	s.data[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.data, v)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.data[v]
	return ok
}

func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.data {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *Set[T]) Items() []T {
	res := make([]T, 0, len(s.data))
	for k := range s.data {
		res = append(res, k)
	}
	return res
}

func NewSet[T comparable](initial ...T) *Set[T] {
	res := &Set[T]{
		data: map[T]struct{}{},
	}
	for _, v := range initial {
		res.data[v] = struct{}{}
	}
	return res
}
