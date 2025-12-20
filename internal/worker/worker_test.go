package worker

import (
	mock "Gonoty/internal/queue/test_mock"
	"context"
	"testing"
)

func TestWorker1(t *testing.T) {
	ctx := context.Background()

	m := mock.NewMockStorage(1, 1000)
	task, _ := m.Dequeue(ctx)

	w := New(m)
	w.processTask(ctx, task)
}

func TestWorker10x100(t *testing.T) {
	ctx := context.Background()

	m := mock.NewMockStorage(10, 100)
	tasks, _ := m.DequeueBatch(ctx, 10)

	w := New(m)
	w.processBatch(ctx, tasks)
}
