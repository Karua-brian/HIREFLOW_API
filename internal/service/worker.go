package service

import (
	"context"
	"fmt"
	"job_board/internal/domain"
	"sync"
	"time"
)

type JobWorker interface {
	Enqueue(ev domain.ApplicationEvent)
}

// wokrker handles background application events
type Worker struct {
	queue chan domain.ApplicationEvent
	workers    int
	wg 		   sync.WaitGroup	
}

// Create a new worker with buffer size
func NewWorker(bufferSize int, workers int) *Worker {
	return &Worker{
		queue: make(chan domain.ApplicationEvent, bufferSize),
		workers: workers,
	}
}

// Start launches background goroutine
func (w *Worker) Start() {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)

		go func (workerID int)  {
			defer w.wg.Done()
			
			for ev := range w.queue {
				// Simulate processing
				fmt.Printf("Worker %d processing job %d\n", workerID, ev.JobID)
			}
		}(i)
	
	}
}

// Enqueue adds a task to the queue
func (w *Worker) Enqueue(ev domain.ApplicationEvent) {
	w.queue <- ev
}

// Stop closes background goroutine to exit
func (w *Worker) Stop() {
	close(w.queue) // closing channel signals goroutine to exit
	w.wg.Wait()
}

func StartTokenCleanUp(ctx context.Context, authService AuthService) {

	// Start a ticker to run cleanup every hour
	ticker := time.NewTicker(1 * time.Hour) 

	go func ()  {
		for {
			select {
			case <-ticker.C:
				err := authService.CleanupExpiredTokens(ctx)
				if err != nil {
					fmt.Printf("Error cleaning up expired tokens: %v\n", err)
				} else {
					fmt.Println("Expired tokens cleaned up successfully")
				}
			case <-ctx.Done():
				fmt.Println("Stopping token cleanup worker")
				ticker.Stop()
				return
			}
		}
	}()
}