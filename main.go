package main

import "fmt"

func producer(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	fmt.Printf("producer done")
}

func consumer(ch <-chan int) {
	for num := range ch {
		fmt.Println(num)
	}
}

func main() {
	ch := make(chan int, 5)
	fmt.Printf("producer starts sending integer to channel")
	go producer(ch)
}
