package sender

import (
	"Gonoty/internal/models"
	"context"
	"fmt"
	"net/smtp"
)

type Sender interface {
	Send() error
}

type SMTPSender struct {
	host      string
	port      int
	fromEmail string
	auth      smtp.Auth
}

type SMTPSenderConfig struct {
	Host      string
	Port      int
	FromEmail string
	Auth      smtp.Auth
}

func NewSender(conf SMTPSenderConfig) *SMTPSender {
	return &SMTPSender{
		host:      conf.Host,
		port:      conf.Port,
		fromEmail: conf.FromEmail,
		auth:      conf.Auth,
	}
}

func (s *Sender) Send(ctx context.Context, tasksChan <-chan []models.Task) {
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
