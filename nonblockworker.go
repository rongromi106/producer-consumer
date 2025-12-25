package main

import (
	"fmt"
	"sync"
	"time"
)

func nonblockworker(ch <-chan int, id int) {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Printf("worker %d channel closed, return\n", id)
				return
			}
			fmt.Printf("worker %d received %d\n", id, v)

			// Reset the idle timer because we just did work.
			if !timer.Stop() {
				// Drain the timer channel if it already fired.
				select {
				case <-timer.C:
				// this is not busy wait as we did not put it in a forloop
				default:
					// can have deadlock when timer is stopped but value is not set into timer.C so it will not be drained
					// but timer is also stopped
				}
			}
			timer.Reset(500 * time.Millisecond)

		case <-timer.C:
			fmt.Printf("worker %d has no work for 500ms, return\n", id)
			return
		}
	}
}

func NonBlockWorkerDemo() {
	ch := make(chan int, 15)
	var wg sync.WaitGroup
	workerCount := 4
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wg.Done()
			nonblockworker(ch, id)
		}()
	}
	go producer(ch, 1, 10, 0, true)
	wg.Wait()
}
