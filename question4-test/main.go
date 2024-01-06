package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	buffer      []byte
	bufferMutex sync.RWMutex
)

func readData(id int) {
	for {
		bufferMutex.RLock()
		fmt.Printf("Reader %d reading data: %v\n", id, buffer)
		bufferMutex.RUnlock()
		time.Sleep(1 * time.Second) // Simulate some reading processing time
	}
}

func writeData(id int, data byte) {
	for {
		bufferMutex.Lock()
		buffer = append(buffer, data)
		fmt.Printf("Writer %d writing data: %v\n", id, buffer)
		bufferMutex.Unlock()
		time.Sleep(2 * time.Second) // Simulate some writing processing time
	}
}

func main() {
	for i := 1; i <= 8; i++ {
		go readData(i)
	}

	for i := 1; i <= 2; i++ {
		go writeData(i, byte(i))
	}

	select {} // Block forever, to keep goroutines running
}
