package main

import (
	"Gonoty/internal/handler"
	"Gonoty/internal/queue"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// var redisConfig = redisstorage.RedisConfig{
// 	Addr:     "localhost:6379",
// 	Password: "",
// 	DB:       0,
// }

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	q, err := queue.New(ctx, queue.Redis)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer q.Close()

	// Handlers
	taskHandler := handler.NewTaskHandler(ctx, q)
	r.Post("/send", taskHandler.AddToQueue)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gonoty"))
	})

	// to env load with Getenv()!
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	// q.Client.FlushDB(ctx) // !!! delete this !!!

	go func() {
		log.Println("Server starting on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
