package benchmark

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"protobench/internal/model"
)

type Runner struct {
	messageCount int
	clients      map[string]model.Protocol
}

func NewRunner(messageCount int) *Runner {
	return &Runner{
		messageCount: messageCount,
		clients:     make(map[string]model.Protocol),
	}
}

func (r *Runner) AddProtocol(name string, protocol model.Protocol) {
	r.clients[name] = protocol
}

func generateTestMessage(id int) *model.Message {
	// Create a realistic payload with ~100KB of data
	content := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		content = append(content, fmt.Sprintf(
			"field%d: This is a detailed data record with multiple fields that might represent a database row or event log. "+
			"Including various data types and lengths to simulate real application data. Current iteration: %d, "+
			"Additional padding to reach desired size with some random values: %d-%d-%d",
			i, i, i*2, i*3, i*4))
	}

	return &model.Message{
		ID:        fmt.Sprintf("msg-%d", id),
		Timestamp: time.Now(),
		Content:   strings.Join(content, "\n"),
		Number:    int64(id),
		IsValid:   true,
	}
}

func (r *Runner) RunBenchmark() []Result {
	var results []Result
	
	for name, protocol := range r.clients {
		start := time.Now()
		errors := 0
		received := make(map[int]bool)

		for i := 0; i < r.messageCount; i++ {
			msg := generateTestMessage(i)
			if err := protocol.SendMessage(msg); err != nil {
				errors++
			} else {
				received[i] = true
			}
		}

		duration := time.Since(start)
		messagesPerSecond := float64(r.messageCount) / duration.Seconds()

		// Check for missing messages
		missing := 0
		for i := 0; i < r.messageCount; i++ {
			if !received[i] {
				missing++
			}
		}

		results = append(results, Result{
			Protocol:         name,
			TotalTime:       duration,
			MessagesPerSecond: messagesPerSecond,
			Errors:          errors,
			Missing:         missing,
		})
	}

	return results
}

func (r *Runner) benchmarkProtocol(name string, protocol model.Protocol) Result {
	var wg sync.WaitGroup
	errorCount := 0
	var errorCountMutex sync.Mutex

	start := time.Now()

	// Create channel for concurrent message sending
	msgChan := make(chan struct{}, r.messageCount)

	// Launch workers
	for i := 0; i < r.messageCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			msg := &model.Message{
				ID:        fmt.Sprintf("msg-%d", id),
				Timestamp: time.Now(),
				Content:   "Hello, World!",
				Number:    int64(id),
				IsValid:   true,
			}

			if err := protocol.SendMessage(msg); err != nil {
				errorCountMutex.Lock()
				errorCount++
				errorCountMutex.Unlock()
			}
			msgChan <- struct{}{}
		}(i)
	}

	wg.Wait()
	close(msgChan)

	duration := time.Since(start)
	messagesPerSecond := float64(r.messageCount) / duration.Seconds()

	return Result{
		Protocol:         name,
		TotalTime:       duration,
		MessagesPerSecond: messagesPerSecond,
		Errors:          errorCount,
		Missing:         0,
	}
} 
