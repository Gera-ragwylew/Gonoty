package main

import (
	"Gonoty/internal/queue"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// var MailHog = sender.SMTPSenderConfig{
// 	Host:      "localhost",
// 	Port:      1025,
// 	FromEmail: "test@example.com",
// 	Auth:      nil,
// }

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q, err := queue.New("redis")
	if err != nil {
		fmt.Println(err)
		return
	}

	ProcessTasks(ctx, q)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func ProcessTasks(ctx context.Context, storage queue.Queue) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task processor stopped")
			return

		default:
			// 3. BRPop из очереди
			task, err := storage.Dequeue(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, v := range task.Recipients {
				// 5. Отправка через net/smtp
				//auth := smtp.PlainAuth("", task.SMTPUser, task.SMTPPass, task.SMTPHost)
				msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", task.Subject, task.Body)
				err = smtp.SendMail(fmt.Sprintf("%s:%d", "localhost", 1025), nil, task.FromEmail, []string{v.Email}, []byte(msg)) // to mailhog

				if err != nil {
					log.Printf("Send failed: %v", err)
				} else {
					log.Printf("Email sent to %s", v.Email)
				}
			}
		}
	}
}
