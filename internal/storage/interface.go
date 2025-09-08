package storage

import (
	"context"

	"Gonoty/internal/models"
)

// Storage интерфейс для абстракции хранилища
type Storage interface {
	// GetPendingTasks возвращает задачи со статусом pending с лимитом
	GetPendingTasks(ctx context.Context, limit int) ([]models.Task, error)

	// UpdateTaskStatus обновляет статус задачи
	UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) error

	// UpdateTasksStatusBatch массовое обновление статусов
	UpdateTasksStatusBatch(ctx context.Context, taskIDs []string, status models.TaskStatus) error
}
