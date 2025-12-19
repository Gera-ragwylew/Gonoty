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
	w.Start(ctx)

	<-ctx.Done()
	// w.Shoutdown()
	log.Println("Server stopped")
}
