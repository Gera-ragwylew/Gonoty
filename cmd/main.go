package main

import (
	"Gonoty/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	// Создаем роутер Chi
	r := chi.NewRouter()

	// Подключаем полезные middleware (логирование, сжатие и т.д.)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Регистрируем ОЧЕНЬ точные пути
	r.Post("/api/send-email", handler.SendEmailHandler) // ТОЛЬКО POST /api/send-email

	// Обработчик для корня "/"
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gonoty API"))
	})

	// Обработчик для всего остального (404 Not Found)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Запускаем сервер
	http.ListenAndServe(":8080", r)
}
