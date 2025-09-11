package main

import (
	"Gonoty/internal/storage/redisstorage"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &redisstorage.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	storage, err := redisstorage.NewRedisStorage(*config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// task := models.Task{
	// 	ID:      fmt.Sprintf("task-%d", 1488),
	// 	Subject: fmt.Sprintf("Test Subject %d", 1),
	// 	Body: models.EmailBody{
	// 		Text: fmt.Sprintf("Email body %d", 1),
	// 	},
	// 	Recipients: []models.Recipient{
	// 		{Email: fmt.Sprintf("user%d@example.com", 1)},
	// 		{Email: fmt.Sprintf("test%d@example.com", 1)},
	// 	},
	// 	FromEmail: "noreply@myapp.com",
	// 	Status:    models.StatusPending,
	// 	CreatedAt: time.Now().Add(-time.Duration(1) * time.Hour),
	// 	UpdatedAt: time.Now().Add(-time.Duration(1) * time.Hour),
	// }

	// storage.Create(ctx, task)
	// storage.List(ctx)
	// storage.Delete(ctx, "1488")
	storage.List(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
	time.Sleep(100 * time.Millisecond)
}
