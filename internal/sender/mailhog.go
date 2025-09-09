package sender

import (
	"Gonoty/internal/models"
	"context"
	"fmt"
	"net/smtp"
)

type Sender struct {
	Host      string
	Port      int
	FromEmail string
	auth      smtp.Auth
}

type Config struct {
	Host      string
	Port      int
	FromEmail string
	Auth      smtp.Auth
}

func NewSender(conf Config) *Sender {
	return &Sender{
		Host:      conf.Host,
		Port:      conf.Port,
		FromEmail: conf.FromEmail,
		auth:      conf.Auth,
	}
}

func (s *Sender) ProcessTasks(ctx context.Context, tasksChan <-chan []models.Task) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task processor stopped")
			return

		case tasks := <-tasksChan:
			fmt.Printf("Processor: Received %d tasks for processing\n", len(tasks))
			for _, task := range tasks {
				fmt.Printf("  - Task %s: %d recipients\n", task.ID, len(task.Recipients))
				to := make([]string, 0)
				for _, v := range task.Recipients {
					to = append(to, v.Email)
				}
				smtp.SendMail(
					fmt.Sprintf("%s:%d", s.Host, s.Port),
					s.auth,
					s.FromEmail,
					to,
					[]byte(task.Body.Text),
				)
			}
		}
	}
}
