package main

import (
	"context"
	"errors"
	"sync"
)

type BlockingQueue[T any] struct {
	capacity int
	mu       sync.Mutex
	notempty sync.Cond
	notfull  sync.Cond
	// using a slice wastes memory by allocation / de-alloc
	// but this is good enough for our use case for now
	arr []T
}

func NewBlockingQueue[T any](cap int) *BlockingQueue[T] {
	return &BlockingQueue[T]{
		capacity: cap,
		arr:      make([]T, cap),
	}
}

func (q *BlockingQueue[T]) Put(ctx context.Context, v T) error {
	q.mu.Lock()
	for len(q.arr)+1 > q.capacity {
		q.notfull.Wait()
		select {
		case <-ctx.Done():
			q.mu.Unlock()
			return errors.New("context is canceled")
		}
	}
	q.arr = append(q.arr, v)
	q.notempty.Signal()
	q.mu.Unlock()
	return nil
}

func (q *BlockingQueue[T]) Take(ctx context.Context) (T, error) {
	q.mu.Lock()
	for len(q.arr) == 0 {
		q.notempty.Wait()
		select {
		case <-ctx.Done():
			q.mu.Unlock()
			var zero T
			return zero, errors.New("context is canceled")
		}
	}
	top := q.arr[0]
	q.arr = q.arr[1:]
	q.notfull.Signal()
	q.mu.Unlock()
	return top, nil
}

func (q *BlockingQueue[T]) Size() int {
	current_size := 0
	q.mu.Lock()
	current_size = len(q.arr)
	q.mu.Unlock()
	return current_size
}
