package storage

import (
	"context"

	"Gonoty/internal/models"
)

type Storage interface {
	Add(ctx context.Context, task models.Task) error
}

// type Storage interface {
// 	GetPendingTasks(ctx context.Context, limit int) ([]models.Task, error)
// 	UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) error
// 	UpdateTasksStatusBatch(ctx context.Context, taskIDs []string, status models.TaskStatus) error
// }
