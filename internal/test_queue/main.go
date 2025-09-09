package main

import (
	"Gonoty/internal/models"
	"Gonoty/internal/scouter"
	"Gonoty/internal/sender"
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

	MailHogConfid := sender.Config{
		Host:      "localhost",
		Port:      1025,
		FromEmail: "test@example.com",
		Auth:      nil,
	}

	scout := scouter.NewScouter(mockStorage, scouterConfig)
	sender := sender.NewSender(MailHogConfid)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scout.Start(ctx)

	go sender.ProcessTasks(ctx, tasksChan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
	time.Sleep(100 * time.Millisecond)
}
