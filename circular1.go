package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	// Create a buffered channel with a size of 1 to act as the circular queue
	queue := make(chan int, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go consumer(queue, &wg, start)

	go producer(queue)

	wg.Wait()

	end := time.Now()

	timeTaken := end.Sub(start)

	fmt.Printf("Time taken : %s\n", timeTaken)
}

func producer(queue chan<- int) {
	// for i := 0; ; i++ {
	// 	// Push data into the circular queue
	// 	queue <- i

	// 	fmt.Printf("Produced: %d\n", i)
	// }
	queue <- 10

}

func consumer(queue <-chan int, wg *sync.WaitGroup, start time.Time) {
	defer wg.Done()

	// for {
	data := <-queue
	elapsed := time.Since(start)

	fmt.Printf("Consumed: %d, Elapsed Time: %s\n", data, elapsed)
	// }
}

// Avg Time Taken to push 1 element - Consumed: 10, Elapsed Time: 37.167µs
// Time taken : 240.375µs
