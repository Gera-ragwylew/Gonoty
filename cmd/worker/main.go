package main

import (
	"Gonoty/internal/queue"
	"Gonoty/internal/worker"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
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

	w := worker.New(ctx, q)
	w.Start()

	<-ctx.Done()
	log.Println("Shutting down...")
	time.Sleep(time.Duration(time.Second * 2))
	log.Println("Server stopped")
}
