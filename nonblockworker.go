package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func nonblockproducer(ch chan<- int, start int, end int, id int) {
	rand.Seed(time.Now().UnixNano())
	for i := start; i <= end; i++ {
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Printf("producer %d sleeping %v before sending %d\n", id, delay, i)
		time.Sleep(delay)
		select {
		case ch <- i:
			fmt.Printf("Producer %d sent %d completed\n", id, i)
		default:
			// 	•	这个 i 永远不会被处理（没有重试), expect dropped message
			fmt.Printf("producer failed to send %d\n", i)
		}
	}
}

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
	ch := make(chan int, 5)
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
	var wgP sync.WaitGroup
	wgP.Add(1)
	go func() {
		defer wgP.Done()
		nonblockproducer(ch, 1, 10, 0)
	}()
	wgP.Wait()
	close(ch)

	wg.Wait()
}
