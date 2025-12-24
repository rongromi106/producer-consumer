package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	fmt.Printf("producer done\n")
	close(ch)
}

func consumer(ch <-chan int, id int) {
	/*
		range ch 这样写会导致consumer一直等待ch close才会退出
		不然就一直阻塞
		所以要在producer中close(ch)
	*/

	for num := range ch {
		fmt.Printf("Consumer id: %d, number: %d\n", id, num)
	}

}

func main() {
	ch := make(chan int, 5)
	fmt.Printf("producer starts sending integer to channel\n")
	var wg sync.WaitGroup
	consumerCount := 2
	wg.Add(consumerCount)

	for i := 0; i < consumerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wg.Done()
			consumer(ch, id)
		}()
	}
	go producer(ch)
	wg.Wait()
	fmt.Printf("main go routine done\n")
}
