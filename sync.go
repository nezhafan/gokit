package gokit

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// 并发执行 tasks 中的每个任务，最多同时运行 limit 个。
func RunAll[T any](ctx context.Context, limit int64, tasks []T, fn func(context.Context, T)) {
	if len(tasks) == 0 || limit <= 0 {
		return
	}

	var wg sync.WaitGroup
	var sem = semaphore.NewWeighted(limit)

	for _, task := range tasks {
		if err := sem.Acquire(ctx, 1); err != nil {
			return
		}

		wg.Add(1)
		go func(task T) {
			defer wg.Done()
			defer sem.Release(1)
			fn(ctx, task)
		}(task)
	}

	wg.Wait()
}

// 复用池
type PoolValue[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

func NewPool[T any](reset func(*T)) *PoolValue[T] {
	return &PoolValue[T]{
		pool:  sync.Pool{New: func() any { return new(T) }},
		reset: reset,
	}
}

func (p *PoolValue[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *PoolValue[T]) Put(v *T) {
	if p.reset != nil {
		p.reset(v)
	}
	p.pool.Put(v)
}

// 单例值
type OnceValue[T any] struct {
	once sync.Once
	init func() T
	val  T
}

func NewOnceValue[T any](fn func() T) *OnceValue[T] {
	return &OnceValue[T]{init: fn}
}

func (o *OnceValue[T]) Get() T {
	o.once.Do(func() {
		if o.init != nil {
			o.val = o.init()
		}
	})
	return o.val
}
