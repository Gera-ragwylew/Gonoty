package main

import (
	"Gonoty/internal/models"
	"Gonoty/internal/queue"
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os/signal"
	"sync"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	q, err := queue.New(ctx, queue.Redis)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer q.Close()

	go func() {
		ProcessTasks(ctx, q)
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	log.Println("Server stopped")
}

func ProcessTasks(ctx context.Context, storage queue.Queue) {
	sem := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			fmt.Println("Task processor stopped")
			return

		default:
			task, err := storage.Dequeue(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				log.Printf("Dequeue error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			wg.Add(1)
			go ProcessTask(ctx, wg, sem, task)
		}
	}
}

func ProcessTask(ctx context.Context, wg *sync.WaitGroup, sem chan struct{}, task models.Task) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	results := make(chan error, len(task.Recipients))
	var taskWg sync.WaitGroup

	for _, recipient := range task.Recipients {
		taskWg.Add(1)
		go func(r models.Recipient) {
			defer taskWg.Done()

			select {
			case <-taskCtx.Done():
				results <- taskCtx.Err()
			default:
				err := sendEmail(taskCtx, task, r)
				results <- err
			}
		}(recipient)
	}

	go func() {
		taskWg.Wait()
		close(results)
	}()

	var successCount, failCount int
	for err := range results {
		if err != nil {
			failCount++
			log.Printf("Send failed for task %s: %v", task.ID, err)
		} else {
			successCount++
		}
	}

	log.Printf("Task %s completed: %d success, %d failed", task.ID, successCount, failCount)

	// if failCount == 0 {
	// 	storage.MarkAsDone(task.ID)
	// } else if successCount > 0 {
	// 	storage.MarkAsFailed(task.ID, fmt.Sprintf("partial failure: %d/%d succeeded", successCount, len(task.Recipients)))
	// } else {
	// 	storage.MarkAsFailed(task.ID, "all sends failed")
	// }
}

func sendEmail(ctx context.Context, task models.Task, recipient models.Recipient) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result := make(chan error, 1)

	go func() {
		msg := fmt.Sprintf("Subject: %s\r\nTo: %s\r\n\r\n%s",
			task.Subject, recipient.Email, task.Body)

		err := smtp.SendMail(
			fmt.Sprintf("%s:%d", "mailhog", 1025),
			nil,
			task.FromEmail,
			[]string{recipient.Email},
			[]byte(msg),
		)
		result <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("send timeout for %s", recipient.Email)
	case err := <-result:
		if err != nil {
			return fmt.Errorf("send to %s failed: %w", recipient.Email, err)
		}
		log.Printf("Email sent to %s", recipient.Email)
		return nil
	}
}
