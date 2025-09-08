package scouter

import (
	"context"
	"fmt"
	"time"

	"Gonoty/internal/models"
	"Gonoty/internal/storage"
)

// Scouter сервис для поиска pending задач
type Scouter struct {
	storage    storage.Storage
	batchSize  int
	interval   time.Duration
	outputChan chan<- []models.Task // Канал для найденных задач
}

// Config конфигурация скаутера
type Config struct {
	BatchSize  int
	Interval   time.Duration
	OutputChan chan<- []models.Task
}

func NewScouter(strg storage.Storage, cfg Config) *Scouter {
	return &Scouter{
		storage:    strg,
		batchSize:  cfg.BatchSize,
		interval:   cfg.Interval,
		outputChan: cfg.OutputChan,
	}
}

// Start запускает скаутер
func (s *Scouter) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	fmt.Printf("Scouter started. Checking every %v with batch size %d\n",
		s.interval, s.batchSize)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Scouter stopped")
			return

		case <-ticker.C:
			s.scout(ctx)
		}
	}
}

// scout ищет pending задачи и отправляет в output канал
func (s *Scouter) scout(ctx context.Context) {
	// 1. Получаем pending задачи
	tasks, err := s.storage.GetPendingTasks(ctx, s.batchSize)
	if err != nil {
		fmt.Printf("Scouter error: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("Scouter: No pending tasks found")
		return
	}

	fmt.Printf("Scouter: Found %d pending tasks\n", len(tasks))

	// 2. Меняем статус на processing
	taskIDs := make([]string, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	if err := s.storage.UpdateTasksStatusBatch(ctx, taskIDs, models.StatusProcessed); err != nil {
		fmt.Printf("Scouter error updating status: %v\n", err)
		return
	}

	// 3. Отправляем задачи в канал для обработки
	select {
	case s.outputChan <- tasks:
		fmt.Printf("Scouter: Sent %d tasks to processor\n", len(tasks))
	case <-ctx.Done():
		fmt.Println("Scouter: Context cancelled while sending tasks")
	case <-time.After(5 * time.Second):
		fmt.Println("Scouter: Timeout sending tasks to processor")
	}
}

// GetStats возвращает статистику (для мониторинга)
func (s *Scouter) GetStats(ctx context.Context) (int, error) {
	tasks, err := s.storage.GetPendingTasks(ctx, 1000) // Большой лимит для подсчета
	if err != nil {
		return 0, err
	}
	return len(tasks), nil
}
