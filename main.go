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

func consumer(ch <-chan int) {
	/*
		range ch 这样写会导致consumer一直等待ch close才会退出
		不然就一直阻塞
		所以要在producer中close(ch)
	*/
	for num := range ch {
		fmt.Println(num)
	}

}

func main() {
	ch := make(chan int, 5)
	fmt.Printf("producer starts sending integer to channel\n")
	go producer(ch)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		consumer(ch)
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("main go routine done\n")
}
