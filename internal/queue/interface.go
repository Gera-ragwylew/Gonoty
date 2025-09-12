package queue

import (
	"context"
	"fmt"

	"Gonoty/internal/models"
	"Gonoty/internal/queue/redisstorage"
)

type Queue interface {
	Enqueue(ctx context.Context, task models.Task) error
	Dequeue(ctx context.Context) (models.Task, error)
	// MarkAsDone(ctx context.Context, task models.Task) error
}

type Type string

const (
	Redis       Type = "redis"
	Postgres    Type = "postgres"
	Mock        Type = "memory"
	FileStorage Type = "file storage"
)

func New(name Type) (Queue, error) {
	switch name {
	case Redis:
		return redisstorage.NewRedisStorage()
	default:
		return nil, fmt.Errorf("unsupported queue type: %s", name)
	}

}
