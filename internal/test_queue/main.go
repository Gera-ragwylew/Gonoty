package main

import (
	"Gonoty/internal/models"
	"Gonoty/internal/scouter"
	"Gonoty/internal/storage"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	mockStorage, _ := storage.NewStorage(storage.Mock)

	tasksChan := make(chan []models.Task, 10)

	scouterConfig := scouter.Config{
		BatchSize:  50,
		Interval:   10 * time.Second,
		OutputChan: tasksChan,
	}

	scout := scouter.NewScouter(mockStorage, scouterConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scout.Start(ctx)

	go processTasks(ctx, tasksChan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func processTasks(ctx context.Context, tasksChan <-chan []models.Task) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task processor stopped")
			return

		case tasks := <-tasksChan:
			fmt.Printf("Processor: Received %d tasks for processing\n", len(tasks))
			time.Sleep(time.Second * 30) // job
			for _, task := range tasks {
				fmt.Printf("  - Task %s: %d recipients\n", task.ID, len(task.Recipients))
			}
		}
	}
}
