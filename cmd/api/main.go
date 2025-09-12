package main

import (
	"Gonoty/internal/handler"
	"Gonoty/internal/queue"
	"Gonoty/internal/queue/redisstorage"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var redisConfig = redisstorage.RedisConfig{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	q, err := queue.New("redis")
	if err != nil {
		fmt.Println(err)
		return
	}
	// q.Client.FlushDB(ctx) // !!! delete this !!!

	taskHandler := handler.NewTaskHandler(q)

	r.Post("/send", taskHandler.SendEmail)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gonoty"))
	})

	// Запускаем сервер
	http.ListenAndServe(":8080", r)
}
