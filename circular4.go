package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	QueueCount    = 10000 // Number of circular queues
	QueueSize     = 10    // Size of each circular queue
	ConsumerCount = 10000 // Number of consumers
	DataCount     = 10000 // Total data items to produce
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup

	queues := make([]chan int, QueueCount)
	for i := 0; i < QueueCount; i++ {
		queues[i] = make(chan int, QueueSize)
		wg.Add(1)
		go producer(queues[i], &wg, i)
	}

	// Create and manage consumers
	for i := 0; i < ConsumerCount; i++ {
		wg.Add(1)
		go consumer(queues[i%QueueCount], &wg, start, DataCount/ConsumerCount)
	}

	// Wait for all producers and consumers to finish
	wg.Wait()

	end := time.Now()

	timeTaken := end.Sub(start)

	fmt.Printf("Total execution time : %s\n", timeTaken)
}

func producer(queue chan<- int, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	startId := id * (DataCount / QueueCount)
	endId := (id + 1) * (DataCount / QueueCount)

	for i := startId; i < endId; i++ {
		queue <- i
		fmt.Printf("Producer %d: Produced: %d\n", id, i)
	}
}

func consumer(queue <-chan int, wg *sync.WaitGroup, start time.Time, itemCount int) {
	defer wg.Done()

	for i := 0; i < itemCount; i++ {
		data := <-queue
		fmt.Printf("Consumed: %d,", data)
	}

}

// Total execution time : 278.428875ms
