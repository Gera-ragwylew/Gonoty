package main

import (
	"Gonoty/internal/handler"
	"Gonoty/internal/models"
	"Gonoty/internal/sender"
	"Gonoty/internal/storage/redisstorage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var redisConfig = redisstorage.RedisConfig{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

var MailHogConfid = sender.Config{
	Host:      "localhost",
	Port:      1025,
	FromEmail: "test@example.com",
	Auth:      nil,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	storage, err := redisstorage.NewRedisStorage(redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	storage.Client.FlushDB(ctx)

	taskHandler := handler.NewTaskHandler(storage)

	go func() {
		for {
			// 3. BRPop из очереди
			result, err := storage.Client.BRPop(context.Background(), 0, "email_queue").Result()
			if err != nil {
				log.Printf("Redis error: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// 4. Парсинг задачи
			var task models.Task
			json.Unmarshal([]byte(result[1]), &task)
			for _, v := range task.Recipients {
				// 5. Отправка через net/smtp
				//auth := smtp.PlainAuth("", task.SMTPUser, task.SMTPPass, task.SMTPHost)
				msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", task.Subject, task.Body)
				err = smtp.SendMail("localhost:1025", nil, task.FromEmail, []string{v.Email}, []byte(msg)) //mailhog

				if err != nil {
					log.Printf("Send failed: %v", err)
				} else {
					log.Printf("Email sent to %s", v.Email)
				}
			}

		}
	}()

	// Регистрируем ОЧЕНЬ точные пути
	r.Post("/send", taskHandler.SendEmail)

	// Обработчик для корня "/"
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gonoty"))
	})

	// Запускаем сервер
	http.ListenAndServe(":8080", r)
}
