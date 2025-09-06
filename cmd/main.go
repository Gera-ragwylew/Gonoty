package main

import (
	"Gonoty/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	taskHandler := handler.NewTaskHandler()

	// Регистрируем ОЧЕНЬ точные пути
	r.Post("/send", taskHandler.SendEmail)

	// Обработчик для корня "/"
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gonoty"))
	})

	// Запускаем сервер
	http.ListenAndServe(":8080", r)
}
