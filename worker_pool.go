package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrNotStarted = errors.New("workerpool: not started")
	ErrClosed     = errors.New("workerpool: closed")
	ErrQueueFull  = errors.New("workerpool: queue full") // 这版默认阻塞，不会用到；留着扩展用
	ErrNilTask    = errors.New("workerpool: nil task")
)

type Task func(ctx context.Context) error

type Pool struct {
	taskCh  chan Task
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	started atomic.Bool
	closed  atomic.Bool
	workers int
}

func NewWorkerPool(workers, queueSize int) (*Pool, error) {
	if workers <= 0 {
		return nil, errors.New("workerpool: workers must be > 0")
	}
	if queueSize < 0 {
		return nil, errors.New("workerpool: queueSize must be >= 0")
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		taskCh:  make(chan Task, queueSize),
		ctx:     ctx,
		cancel:  cancel,
		workers: workers,
	}, nil
}

func (p *Pool) Start() {
	if !p.started.CompareAndSwap(false, true) {
		return
	}
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.workerLoop()
	}
}

// 默认队列满了塞不进去就阻塞
func (p *Pool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return ErrNilTask
	}
	if !p.started.Load() {
		return ErrNotStarted
	}
	if p.closed.Load() {
		return ErrClosed
	}
	select {
	case p.taskCh <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return ErrClosed
	}
}

func (p *Pool) Shutdown(ctx context.Context) error {
	// 只允许 shutdown 一次：标记 closed + cancel + close(taskCh)
	if p.closed.CompareAndSwap(false, true) {
		// 先 cancel 广播 worker 也能退出（但我们是 draining 模式，下面会 close(taskCh) 让 worker 最终退出）
		p.cancel()
		close(p.taskCh)
	}

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (p *Pool) workerLoop() {
	defer p.wg.Done()
	for task := range p.taskCh {
		safeRun(p.ctx, task)
	}
}

// safeRun: 防止 task panic 把 worker 干掉，造成并发度下降或死锁
func safeRun(ctx context.Context, task Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 这版先不把 panic 传播出去，只保证 pool 不崩
			err = errors.New("workerpool: task panicked")
		}
	}()
	return task(ctx)
}
