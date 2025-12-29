package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// simulate pipelining: numbers -> square -> double -> sum

// square worker processes data by squaring it then send to channel
type SquareWorker struct {
	id  int
	in  <-chan int
	out chan<- int
}

func NewSquareWorker(id int, input <-chan int, output chan<- int) *SquareWorker {
	return &SquareWorker{
		id:  id,
		in:  input,
		out: output,
	}
}

func (w *SquareWorker) Work() {
	for num := range w.in {
		fmt.Printf("SquareWorker %d: %d\n", w.id, num*num)
		w.out <- num * num
	}
}

func (w *SquareWorker) WorkWithContext(ctx context.Context) {
	for {
		select {
		case num, ok := <-w.in:
			if !ok {
				fmt.Println("Channel closed.")
				return
				// TODO: error handling
			}
			select {
			case <-ctx.Done():
				return
			case w.out <- num * num:
			}
		case <-ctx.Done():
			return
		}
	}

}

// double worker processes data by doubling it then send to channel
type DoubleWorker struct {
	id  int
	in  <-chan int
	out chan<- int
}

func NewDoubleWorker(id int, input <-chan int, output chan<- int) *DoubleWorker {
	return &DoubleWorker{
		id:  id,
		in:  input,
		out: output,
	}
}

func (w *DoubleWorker) Work() {
	for num := range w.in {
		fmt.Printf("DoubleWorker %d: %d\n", w.id, num*2)
		w.out <- num * 2
	}
}

func (w *DoubleWorker) WorkWithContext(ctx context.Context) {
	for {
		select {
		case num, ok := <-w.in:
			if !ok {
				fmt.Println("Channel closed.")
				return
				// TODO: error handling
			}
			select {
			case <-ctx.Done():
				return
			case w.out <- num * 2:
			}
		case <-ctx.Done():
			return
		}
	}
}

func producer(ch chan<- int, start int, end int, id int, closecchan bool) {
	// 建议只 seed 一次；这里为了 demo 放在函数里也 OK
	rand.Seed(time.Now().UnixNano())

	for i := start; i <= end; i++ {
		// random sleep: 0 ~ 1 seconds
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Printf("producer %d sleeping %v before sending %d\n", id, delay, i)
		time.Sleep(delay)
		ch <- i
	}
	fmt.Printf("producer %d done\n", id)
	if closecchan {
		close(ch)
	}
}

func PipeliningDemo() {
	// numbers -> square -> double -> sum
	numbers := make(chan int, 10)
	squared := make(chan int, 10)
	doubled := make(chan int, 10)

	fmt.Printf("producer starts sending integer to channel\n")
	go producer(numbers, 1, 10, 1, true)

	// Fan-out square workers
	var wgSquare sync.WaitGroup
	squareWorkerCount := 2
	wgSquare.Add(squareWorkerCount)
	for i := 0; i < squareWorkerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wgSquare.Done()
			NewSquareWorker(id, numbers, squared).Work()
		}()
	}
	// Fan-in: close squared when all square workers are done
	go func() {
		wgSquare.Wait()
		close(squared)
	}()

	// Fan-out double workers
	var wgDouble sync.WaitGroup
	doubleWorkerCount := 3
	wgDouble.Add(doubleWorkerCount)
	for i := 0; i < doubleWorkerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wgDouble.Done()
			NewDoubleWorker(id, squared, doubled).Work()
		}()
	}
	// Fan-in: close doubled when all double workers are done
	go func() {
		wgDouble.Wait()
		close(doubled)
	}()

	// Sink: consume doubled concurrently while workers are producing
	sum := 0
	for num := range doubled {
		sum += num
	}
	fmt.Printf("sum: %d\n", sum)
	fmt.Printf("main go routine done\n")
}

func producerWithContext(nums chan int, ctx context.Context, N int) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= N; i++ {
		// random sleep: 0 ~ 1 seconds, this should cause some context cancellation and workers stopping
		delay := time.Duration(rand.Intn(300)) * time.Millisecond
		fmt.Printf("producer sleeping %v before sending %d\n", delay, i)
		time.Sleep(delay)
		select {
		case nums <- i:
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Producer failed to send %d\n", i)
		}
	}
	close(nums)
}

// Day 5: worker with context
func RunPipeline(ctx context.Context, N int, sqaureWorkerCounts int, doubleWorkerCounts int) int {
	// numbers -> square -> double -> sum
	numbers := make(chan int, 10)
	squared := make(chan int, 10)
	doubled := make(chan int, 10)
	fmt.Printf("producer starts sending integer to channel\n")
	go producerWithContext(numbers, ctx, N)

	// Fan-out square workers
	var wgSquare sync.WaitGroup
	wgSquare.Add(sqaureWorkerCounts)
	for i := 0; i < sqaureWorkerCounts; i++ {
		id := i
		go func() {
			defer wgSquare.Done()
			NewSquareWorker(id, numbers, squared).WorkWithContext(ctx)
		}()
	}

	// Fan-in: close squared when all square workers are done
	go func() {
		wgSquare.Wait()
		close(squared)
	}()

	// Fan-out double workers
	var wgDouble sync.WaitGroup
	wgDouble.Add(doubleWorkerCounts)
	for i := 0; i < doubleWorkerCounts; i++ {
		id := i // capture loop variable
		go func() {
			defer wgDouble.Done()
			NewDoubleWorker(id, squared, doubled).WorkWithContext(ctx)
		}()
	}
	// Fan-in: close doubled when all double workers are done
	go func() {
		wgDouble.Wait()
		close(doubled)
	}()

	// Sink: consume doubled concurrently while workers are producing
	sum := 0
	for {
		select {
		case v, ok := <-doubled:
			if !ok {
				fmt.Printf("double channel closed, returning sum normally\n")
				return sum
			}
			sum += v
		case <-ctx.Done():
			fmt.Printf("context canceled, returning partial sum\n")
			return sum
		}
	}
}
