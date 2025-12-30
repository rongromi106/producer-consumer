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
	closed   bool
	// using a slice wastes memory by allocation / de-alloc
	// but this is good enough for our use case for now
	arr []T
}

func NewBlockingQueue[T any](cap int) *BlockingQueue[T] {
	q := &BlockingQueue[T]{
		capacity: cap,
		arr:      make([]T, 0, cap),
	}
	q.notempty.L = &q.mu
	q.notfull.L = &q.mu
	return q
}

// Close marks the queue as closed and wakes all waiters.
// After Close, Put returns an error and Take returns an error once the queue is drained.
func (q *BlockingQueue[T]) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	// Wake both producers and consumers.
	q.notempty.Broadcast()
	q.notfull.Broadcast()
	q.mu.Unlock()
}

func (q *BlockingQueue[T]) Put(ctx context.Context, v T) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if ctx.Err() != nil {
		return ctx.Err()
	}
	if q.closed {
		return errors.New("queue is closed")
	}

	for len(q.arr) >= q.capacity {
		if q.closed {
			return errors.New("queue is closed")
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		q.notfull.Wait()
	}

	if q.closed {
		return errors.New("queue is closed")
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	q.arr = append(q.arr, v)
	q.notempty.Signal()
	return nil
}

func (q *BlockingQueue[T]) Take(ctx context.Context) (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for {
		// Drain takes precedence: if we have data, return it even if ctx is canceled.
		if len(q.arr) > 0 {
			top := q.arr[0]
			q.arr = q.arr[1:]
			q.notfull.Signal()
			return top, nil
		}

		// No buffered items.
		if q.closed {
			var zero T
			return zero, errors.New("queue is closed")
		}

		if err := ctx.Err(); err != nil {
			var zero T
			return zero, err
		}

		// Wait for either data to arrive or the queue to be closed.
		q.notempty.Wait()
	}
}

func (q *BlockingQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.arr)
}
