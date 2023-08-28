package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Stock struct {
	Symbol  string
	Price   float64
	Updated time.Time
}

func main() {
	// Configuration
	maxQueueSize := 10000             // Maximum buffer size for the data queue
	initialQueueSize := 1000          // Initial buffer size
	queueSizeIncreaseThreshold := 0.8 // Queue size increase threshold
	queueSizeDecreaseThreshold := 0.3 // Queue size decrease threshold
	queueSizeIncreaseFactor := 2      // Queue size increase factor
	queueSizeDecreaseFactor := 2      // Queue size decrease factor

	// Create a buffered channel to simulate the data queue.
	dataQueue := make(chan Stock, initialQueueSize)

	// Create a wait group to synchronize goroutines.
	var wg sync.WaitGroup

	// Start the publisher goroutine to continuously push stock data to the queue.
	wg.Add(1)
	go publisher(dataQueue, &wg)

	// Start multiple subscriber goroutines.
	numSubscribers := 5
	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go subscriber(i, dataQueue, &wg)
	}

	// Monitor the queue size and adjust it dynamically.
	wg.Add(1)
	go queueSizeMonitor(dataQueue, maxQueueSize, queueSizeIncreaseThreshold, queueSizeDecreaseThreshold, queueSizeIncreaseFactor, queueSizeDecreaseFactor, &wg)

	// Wait for all goroutines to finish.
	wg.Wait()
}

func publisher(queue chan<- Stock, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// TODO : Simulate getting stock data from NSE (replace with actual data retrieval).
		stockData := generateRandomStockData()
		// Push the stock data to the queue.
		queue <- stockData
		fmt.Printf("Published: %+v\n", stockData)
	}
}

func subscriber(id int, queue <-chan Stock, wg *sync.WaitGroup) {
	defer wg.Done()

	for stockData := range queue {
		fmt.Printf("Subscriber %d received: %+v\n", id, stockData)
		// Just printing the stock data and the subscriber id
	}
}

func generateRandomStockData() Stock {
	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"}
	rand.Seed(time.Now().UnixNano())
	randomSymbol := symbols[rand.Intn(len(symbols))]
	randomPrice := rand.Float64() * 1000
	return Stock{
		Symbol:  randomSymbol,
		Price:   randomPrice,
		Updated: time.Now(),
	}
}

func queueSizeMonitor(queue chan Stock, maxQueueSize int, increaseThreshold, decreaseThreshold float64, increaseFactor, decreaseFactor int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Check the current queue size and calculate the current threshold.
		queueSize := len(queue)
		currentThreshold := float64(queueSize) / float64(maxQueueSize)

		// Dynamically adjust the queue size based on thresholds.
		if currentThreshold > increaseThreshold && maxQueueSize*increaseFactor > queueSize {
			// Increase the queue size and replace the old queue.
			newQueueSize := maxQueueSize * increaseFactor
			fmt.Printf("Increasing queue size to %d\n", newQueueSize)
			queue = adjustQueueSize(queue, newQueueSize)
		} else if currentThreshold < decreaseThreshold && maxQueueSize/decreaseFactor < queueSize {
			// Decrease the queue size and replace the old queue.
			newQueueSize := maxQueueSize / decreaseFactor
			fmt.Printf("Decreasing queue size to %d\n", newQueueSize)
			queue = adjustQueueSize(queue, newQueueSize)
		}

		// Sleep for a while before checking again. Checking the size of queue every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func adjustQueueSize(queue chan Stock, newQueueSize int) chan Stock {
	// Create a new channel with the adjusted size.
	newQueue := make(chan Stock, newQueueSize)

	// Move data from the old queue to the new one in a separate goroutine.
	go func() {
		defer close(newQueue)
		for data := range queue {
			select {
			case newQueue <- data:
			default:
				break
			}
		}
	}()

	return newQueue
}
