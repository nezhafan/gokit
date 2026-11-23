package gokit

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// RunAll 执行 tasks 中的所有任务，最多同时并发 limit 个。
//
// 参数：
//   - ctx：上下文；ctx 取消后停止派发尚未开始的任务（已启动的任务不受影响，会继续跑完）
//   - tasks：任务数据列表，每个元素会作为参数传入 fn
//   - limit：最大并发执行数量
//   - fn：单个任务的执行函数，参数为任务项；fn 内部的错误自行处理，本函数不收集
//
// 返回值：
//   - <-chan bool：带缓冲(1)的 channel，所有任务都被派发并执行完毕
//     收到 true  —— 所有任务均已派发且执行完成（tasks 为空时也为 true）
//     收到 false —— ctx 在派发过程中被取消，部分任务未被派发/执行
//
// 用法：
//
//	allDone := <-RunAll(ctx, tasks, 10, fn) // 此时阻塞
//	if !allDone { /* ctx 被取消，任务未全部执行 */ }
func RunAll[T any](ctx context.Context, tasks []T, limit int64, fn func(task T)) <-chan bool {
	result := make(chan bool, 1)

	// 任务数为0，直接关闭返回
	if len(tasks) == 0 {
		result <- true
		close(result)
		return result
	}

	// 不推荐，但是兼容
	if limit <= 0 {
		limit = int64(len(tasks))
	}

	go func() {
		defer close(result)

		var wg sync.WaitGroup
		sem := semaphore.NewWeighted(limit)
		allDispatched := true

		for _, task := range tasks {
			// 并发信号申请
			if err := sem.Acquire(ctx, 1); err != nil {
				allDispatched = false
				break
			}
			// 总数+1
			wg.Add(1)

			// 并发执行单个任务
			go func(task T) {
				defer wg.Done()      // 总数-1
				defer sem.Release(1) // 并发信号释放
				fn(task)
			}(task)
		}

		wg.Wait()
		result <- allDispatched
	}()

	return result
}

// 复用池
type Pool[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

// 对象复用池 (对 sync.Pool 的简单泛型封装 )
//
//	construct := func() *User { return new(User) } /* 对象构造方法 */
//	reset := func(u *User) { u.Name = "" } /* 对象属性重置 */
//	p := NewPool(construct, reset)
//	u1 := p.Get()
//	defer p.Put(u1)
//	/* do somethine */
func NewPool[T any](construct func() *T, reset func(*T)) *Pool[T] {
	return &Pool[T]{
		pool:  sync.Pool{New: func() any { return construct() }},
		reset: reset,
	}
}

func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *Pool[T]) Put(v *T) {
	if p.reset != nil {
		p.reset(v)
	}
	p.pool.Put(v)
}

// 单例值
type Once[T any] struct {
	once sync.Once
	init func() T
	val  T
}

// 传入初始化方法，该函数懒加载一次。
// RunAll 执行 tasks 中的所有任务，最多同时并发 limit 个。
/*
	u := NewOnce(func() *User { return &User{Time: time.Now().UnixMilli()} })
	fmt.Println(u.Get().Time)
	time.Sleep(time.Millisecond * 10)
	fmt.Println(u.Get().Time)
*/

// 单例懒加载对象 (对 sync.Once 的简单封装 )
//
// 用法：
//
//	db := NewOnce(func() *DB { /* 初始化 */ })
//	db.Get().Select().Where().Find()
func NewOnce[T any](fn func() T) *Once[T] {
	return &Once[T]{init: fn}
}

func (o *Once[T]) Get() T {
	o.once.Do(func() {
		if o.init != nil {
			o.val = o.init()
		}
	})
	return o.val
}

// sync.Map 的泛型封装
type Map[K comparable, V any] struct {
	data sync.Map
}

func (m *Map[K, V]) Set(key K, value V) {
	m.data.Store(key, value)
}

func (m *Map[K, V]) MSet(kv map[K]V) {
	for k, v := range kv {
		m.Set(k, v)
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.data.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (m *Map[K, V]) MGet(keys []K) map[K]V {
	r := make(map[K]V, len(keys))
	for _, k := range keys {
		v, ok := m.data.Load(k)
		if ok {
			r[k] = v.(V)
		}
	}

	return r
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.data.Load(key)
	return ok
}

func (m *Map[K, V]) Del(key K) {
	m.data.Delete(key)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.data.Range(func(key, value any) bool {
		k := key.(K)
		v := value.(V)
		return f(k, v)
	})
}
