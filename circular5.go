package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	QueueCount    = 10000   // Number of circular queues
	QueueSize     = 100     // Size of each circular queue
	ConsumerCount = 10000   // Number of consumers
	DataCount     = 1000000 // Total data items to produce
)

func main() {
	start := time.Now() // Record the start time
	var wg sync.WaitGroup

	queues := make([]chan int, QueueCount)
	for i := 0; i < QueueCount; i++ {
		queues[i] = make(chan int, QueueSize)
		producer(queues[i], &wg, i)
	}

	for i := 0; i < ConsumerCount; i++ {
		go consumerWorker(queues, &wg)
	}
	wg.Wait()
	end := time.Now()

	totalExecutionTime := end.Sub(start)
	fmt.Printf("Total Execution Time: %s\n", totalExecutionTime)
}

func producer(queue chan<- int, wg *sync.WaitGroup, id int) {

	start := id * (DataCount / QueueCount)
	end := (id + 1) * (DataCount / QueueCount)

	for i := start; i < end; i++ {
		queue <- i
		// fmt.Printf("Producer %d: Produced: %d\n", id, i)
	}
}

func consumerWorker(queues []chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for _, queue := range queues {
			data := <-queue
			fmt.Printf("Consumed: %d\n", data)
		}
	}
}

func newFunction(startTime time.Time) time.Duration {
	endTime := time.Since(startTime)
	return endTime
}

//Total Execution Time: 50.703292ms
