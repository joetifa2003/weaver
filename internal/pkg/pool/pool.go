package pool

import "sync"

type Pool[T any] struct {
	sp        sync.Pool
	beforePut func(t T)
}

type Option[T any] func(p *Pool[T])

func WithBeforePut[T any](beforePut func(t T)) Option[T] {
	return func(p *Pool[T]) {
		p.beforePut = beforePut
	}
}

func New[T any](new func() T, options ...Option[T]) *Pool[T] {
	return &Pool[T]{
		sp: sync.Pool{
			New: func() any {
				return new()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.sp.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	if p.beforePut != nil {
		p.beforePut(t)
	}

	p.sp.Put(t)
}
