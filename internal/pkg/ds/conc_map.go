package ds

import (
	"iter"
	"sync"
)

type ConcMap[K comparable, V any] struct {
	m sync.Map
}

func NewConcMap[K comparable, V any]() *ConcMap[K, V] {
	return &ConcMap[K, V]{}
}

func (c *ConcMap[K, V]) Set(k K, v V) {
	c.m.Store(k, v)
}

func (c *ConcMap[K, V]) Get(k K) (V, bool) {
	v, ok := c.m.Load(k)
	if !ok {
		return zero[V](), false
	}
	return v.(V), ok
}

func (c *ConcMap[K, V]) Delete(k K) {
	c.m.Delete(k)
}

func (c *ConcMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.m.Range(func(key, value any) bool {
			return yield(key.(K), value.(V))
		})
	}
}

func zero[T any]() T {
	var zero T
	return zero
}
