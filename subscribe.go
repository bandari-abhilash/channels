package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Stock struct {
	Symbol  string
	Price   float64
	Updated time.Time
}

func main() {
	// Configuration (you can use environment variables or a configuration file)
	queueSize := 100
	numSubscribers := 3
	dataUpdateInterval := 2 * time.Second

	// Create a channel to simulate the data queue.
	dataQueue := make(chan Stock, queueSize)

	// Create a wait group to synchronize goroutines.
	var wg sync.WaitGroup

	// Start the publisher goroutine to continuously push stock data to the queue.
	wg.Add(1)
	go publisher(dataQueue, dataUpdateInterval, &wg)

	// Start multiple subscriber goroutines.
	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go subscriber(i, dataQueue, &wg)
	}

	// Handle graceful shutdown (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // Block until a signal is received

	// Close the dataQueue and wait for all goroutines to finish.
	close(dataQueue)
	wg.Wait()
}

func publisher(queue chan<- Stock, updateInterval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Simulate getting stock data from NSE (replace with actual data retrieval).
		stockData := generateRandomStockData()
		// Push the stock data to the queue.
		select {
		case queue <- stockData:
			fmt.Printf("Published: %+v\n", stockData)
		default:
			fmt.Println("Queue is full. Skipping data.")
		}

		// Sleep for a while to simulate real-time updates.
		time.Sleep(updateInterval)
	}
}

func subscriber(id int, queue <-chan Stock, wg *sync.WaitGroup) {
	defer wg.Done()

	for stockData := range queue {
		fmt.Printf("Subscriber %d received: %+v\n", id, stockData)

		// You can perform your processing here.
		// For example, you can implement logic to buy/sell stocks or analyze data.
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
