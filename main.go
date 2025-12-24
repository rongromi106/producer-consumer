package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int, start int, end int, id int) {
	for i := start; i <= end; i++ {
		ch <- i
	}
	fmt.Printf("producer %d done\n", id)
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
	var wgProd sync.WaitGroup
	consumerCount := 2
	producerCount := 3
	wg.Add(consumerCount)
	wgProd.Add(producerCount)

	for i := 0; i < consumerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wg.Done()
			consumer(ch, id)
		}()
	}

	for i := 0; i < producerCount; i++ {
		id := i // capture loop variable
		go func() {
			defer wgProd.Done()
			producer(ch, id*10, id*10+9, id)
		}()
	}

	wgProd.Wait()
	close(ch)
	wg.Wait()
	fmt.Printf("main go routine done\n")
}
