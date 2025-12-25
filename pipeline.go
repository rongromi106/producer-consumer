package main

import (
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
