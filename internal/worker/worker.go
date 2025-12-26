package worker

import (
	"Gonoty/internal/models"
	"Gonoty/internal/queue"
	"context"
	"fmt"
	"math/rand"
	"mime"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	q               queue.Queue
	config          Config
	targetInFlight  int
	currentInFlight int32
	minBatchSize    int
	maxBatchSize    int
	taskSemaphore   chan struct{}
}

func New(queue queue.Queue) *Worker {
	return &Worker{
		q:               queue,
		config:          NewConfig(),
		targetInFlight:  5,
		currentInFlight: 0,
		minBatchSize:    1,
		maxBatchSize:    5,
		taskSemaphore:   make(chan struct{}, 5),
	}
}

func (w *Worker) Start(ctx context.Context) {
	if err := w.config.LoadFromFile(); err != nil {
		fmt.Println("Config error: %w", err)
	}

	go func() {
		// // Включаем STARTTLS
		// if ok, _ := c.Extension("STARTTLS"); ok {
		// 	if err := c.StartTLS(&tls.Config{
		// 		ServerName: "smtp.yandex.ru",
		// 	}); err != nil {
		// 		fmt.Println("STARTTLS failed: %w", err)
		// 		return
		// 	}
		// }

		// isProcessing := false
		// var sumStart time.Time

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// task, err := w.q.Dequeue(ctx, 1)
				// if err != nil || task.ID == "" {
				// 	if isProcessing {
				// 		isProcessing = false
				// 		fmt.Println("All tasks complete with", time.Since(sumStart))
				// 	}
				// 	fmt.Println(err)
				// 	// add try reconnet
				// 	time.Sleep(100 * time.Millisecond)
				// 	continue
				// }

				// if task.ID != "" && !isProcessing {
				// 	isProcessing = true
				// 	fmt.Println("start sum timer...")
				// 	sumStart = time.Now()
				// }

				// fmt.Println(task.ID, "process...")
				// start := time.Now()

				// auth := smtp.PlainAuth("", "test@yandex.ru", "yandexpsw", "smtp.yandex.ru")
				// if err := c.Auth(auth); err != nil {
				// 	fmt.Println(err)
				// }

				freeSlots := w.targetInFlight - int(atomic.LoadInt32(&w.currentInFlight))
				if freeSlots <= 0 {
					time.Sleep(10 * time.Millisecond)
					continue
				}

				batchSize := min(freeSlots*2, w.maxBatchSize)
				batchSize = max(batchSize, w.minBatchSize)

				tasks, err := w.q.DequeueBatch(ctx, batchSize)
				if err != nil {
					fmt.Println(err)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				if len(tasks) == 0 {
					time.Sleep(50 * time.Millisecond)
					continue
				}

				w.processBatch(ctx, tasks)
			}
		}
	}()

	<-ctx.Done()
	// wg := &sync.WaitGroup{}
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			wg.Wait()
	// 			fmt.Println("Task processor stopped")
	// 			return

	// 		default:
	// 			task, err := w.q.Dequeue(ctx)
	// 			if err != nil {
	// 				if errors.Is(err, context.Canceled) {
	// 					continue
	// 				}
	// 				log.Printf("Dequeue error: %v", err)
	// 				time.Sleep(1 * time.Second)
	// 				continue
	// 			}

	// 			wg.Add(1)
	// 			go w.processTask(ctx, wg, task)
	// 		}
	// 	}
	// }()
}

func (w *Worker) processBatch(ctx context.Context, tasks []models.Task) {
	for _, task := range tasks {
		w.taskSemaphore <- struct{}{}
		atomic.AddInt32(&w.currentInFlight, 1)

		go func() {
			defer func() {
				<-w.taskSemaphore
				atomic.AddInt32(&w.currentInFlight, -1)
			}()

			w.processTask(ctx, task)
		}()
	}
}

func (w *Worker) processTask(ctx context.Context, task models.Task) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)

	// c, err := smtp.Dial("localhost:1025")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer func() {
	// 	err = c.Quit()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	fmt.Println(task.ID, "process...")
	start := time.Now()

	pool := NewSMTPPool(w.config.Address, w.config.MaxConnections)

	for _, r := range task.Recipients {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			c, _ := pool.Get()
			defer pool.Put(c)

			if err := sendEmail(ctx, c, task, r); err != nil {
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("task ", task.ID, "complete with", time.Since(start))

	// defer wg.Done()

	// w.taskSem <- struct{}{}
	// defer func() { <-w.taskSem }()

	// taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	// defer cancel()

	// results := make(chan error, len(task.Recipients))
	// var taskWg sync.WaitGroup

	// for _, recipient := range task.Recipients {
	// 	taskWg.Add(1)
	// 	go func(r models.Recipient) {
	// 		defer taskWg.Done()

	// 		select {
	// 		case <-taskCtx.Done():
	// 			results <- taskCtx.Err()
	// 		default:
	// 			err := sendEmail(taskCtx, task, r)
	// 			results <- err
	// 		}
	// 	}(recipient)
	// }

	// go func() {
	// 	taskWg.Wait()
	// 	close(results)
	// }()

	// var successCount, failCount int
	// for err := range results {
	// 	if err != nil {
	// 		failCount++
	// 		log.Printf("Send failed for task %s: %v", task.ID, err)
	// 	} else {
	// 		successCount++
	// 	}
	// }

	// log.Printf("Task %s completed: %d success, %d failed", task.ID, successCount, failCount)

	// if failCount == 0 {
	// 	storage.MarkAsDone(task.ID)
	// } else if successCount > 0 {
	// 	storage.MarkAsFailed(task.ID, fmt.Sprintf("partial failure: %d/%d succeeded", successCount, len(task.Recipients)))
	// } else {
	// 	storage.MarkAsFailed(task.ID, "all sends failed")
	// }
}

func sendEmail(ctx context.Context, c *smtp.Client, task models.Task, recipient models.Recipient) error {
	// ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()

	// domain := strings.Split(recipient.Email, "@")[1]
	// mxRecords, err := lookupMX(domain)
	// if err != nil {
	// 	return err
	// }

	msg := messageBuilder(task, recipient)

	if err := c.Mail(task.FromEmail); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := c.Rcpt(recipient.Email); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	_, err = wc.Write(msg)
	if err != nil {
		return fmt.Errorf("writing message failed: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("closing data failed: %w", err)
	}

	return nil
	// result <- err
	// }()

	// select {
	// case <-ctx.Done():
	// 	return fmt.Errorf("send timeout for %s", recipient.Email)
	// case err := <-result:
	// 	if err != nil {
	// 		return fmt.Errorf("send to %s failed: %w", recipient.Email, err)
	// 	}
	// 	log.Printf("Email sent to %s", recipient.Email)
	// 	return nil
	// }
}

func messageBuilder(task models.Task, recipient models.Recipient) []byte {
	var msg strings.Builder
	boundary := fmt.Sprintf("boundary-%d", time.Now().UnixNano())

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", task.FromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", recipient.Email))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", task.Subject)))
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString(fmt.Sprintf("Message-ID: <%s@myapp.com>\r\n",
		fmt.Sprintf("%d.%d", time.Now().UnixNano(), rand.Int63())))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	msg.WriteString("\r\n")

	// Plain text
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(task.Body.Text)
	msg.WriteString("\r\n")

	// HTML
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(task.Body.HTML)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return []byte(msg.String())
}

func lookupMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("error looking up MX records: %w", err)
	}

	return mxRecords, nil
}
