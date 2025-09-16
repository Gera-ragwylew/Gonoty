package worker

import (
	"Gonoty/internal/models"
	"Gonoty/internal/queue"
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

const (
	maxConcurrentTasks = 10
)

type Worker struct {
	ctx     context.Context
	q       queue.Queue
	taskSem chan struct{}
	wg      sync.WaitGroup
}

func New(ctx context.Context, queue queue.Queue) *Worker {
	return &Worker{
		ctx:     ctx,
		q:       queue,
		taskSem: make(chan struct{}, maxConcurrentTasks),
	}
}

func (w *Worker) Start() {
	for {
		select {
		case <-w.ctx.Done():
			w.wg.Wait()
			fmt.Println("Task processor stopped")
			return

		default:
			task, err := w.q.Dequeue(w.ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				log.Printf("Dequeue error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			w.wg.Add(1)
			go w.processTask(task)
		}
	}
}

func (w *Worker) processTask(task models.Task) {
	defer w.wg.Done()

	w.taskSem <- struct{}{}
	defer func() { <-w.taskSem }()

	taskCtx, cancel := context.WithTimeout(w.ctx, 5*time.Minute)
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
