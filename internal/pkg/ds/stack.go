package ds

import "iter"

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() T {
	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return res
}

func (s *Stack[T]) Peek() T {
	var z T
	if len(s.data) > 0 {
		z = s.data[len(s.data)-1]
	}
	return z
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Set(i int, v T) {
	s.data[i] = v
}

func (s *Stack[T]) Get(i int) T {
	return s.data[i]
}

func (s *Stack[T]) Iter() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(s.data) - 1; i >= 0; i-- {
			if !yield(i, s.data[i]) {
				return
			}
		}
	}
}

func NewStack[T any](initial ...T) *Stack[T] {
	return &Stack[T]{
		data: initial,
	}
}
