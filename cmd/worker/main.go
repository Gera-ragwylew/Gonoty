package main

import (
	"Gonoty/internal/queue"
	"Gonoty/internal/worker"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	q, err := queue.New(ctx, queue.Redis)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer q.Close()

	w := worker.New(q)
	if err := w.Start(ctx); err != nil {
		fmt.Println(err)
	}

	<-ctx.Done()
	log.Println("Shutting down...")
	log.Println("Server stopped")
}
