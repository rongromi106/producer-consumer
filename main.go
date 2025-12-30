package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// Day 3 content for pipelining
	// PipeliningDemo()

	// Day 4 content
	// Non blocking Worker
	// NonBlockWorkerDemo()

	// Day 5 context
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		res := RunPipeline(ctx, 10, 2, 3)
		fmt.Printf("Results from pipeline run %d\n", res)
	*/

	// Day 6 Blocking Queue
	// Test case 1: single consumer and single producer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	bq := NewBlockingQueue[int](2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		<-ctx.Done()
		bq.Close()
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			fmt.Printf("producer: putting %d\n", i)
			if err := bq.Put(ctx, i); err != nil {
				fmt.Printf("producer error: %v\n", err)
				return
			}
			fmt.Printf("producer: put %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			v, err := bq.Take(ctx)
			if err != nil {
				fmt.Printf("consumer error: %v\n", err)
				return
			}
			fmt.Printf("consumer: got %d\n", v)
			time.Sleep(3 * time.Second)
		}
	}()

	wg.Wait()
	fmt.Printf("Queue size check: %d\n", bq.Size())
}
