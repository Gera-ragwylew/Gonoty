package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"Gonoty/internal/models"
)

type MockStorage struct {
	tasks []models.Task
	mu    sync.RWMutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		tasks: generateMockTasks(100),
	}
}

func (m *MockStorage) GetPendingTasks(ctx context.Context, limit int) ([]models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var pendingTasks []models.Task
	count := 0

	for _, task := range m.tasks {
		if task.Status == models.StatusPending {
			pendingTasks = append(pendingTasks, task)
			count++
			if count >= limit {
				break
			}
		}
	}

	fmt.Printf("MockStorage: Found %d pending tasks (limit: %d)\n", len(pendingTasks), limit)
	return pendingTasks, nil
}

func (m *MockStorage) UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.tasks {
		if m.tasks[i].ID == taskID {
			m.tasks[i].Status = status
			m.tasks[i].UpdatedAt = time.Now()
			fmt.Printf("MockStorage: Updated task %s to status %s\n", taskID, status)
			return nil
		}
	}

	return fmt.Errorf("task not found: %s", taskID)
}

func (m *MockStorage) UpdateTasksStatusBatch(ctx context.Context, taskIDs []string, status models.TaskStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	updated := 0
	for _, taskID := range taskIDs {
		for i := range m.tasks {
			if m.tasks[i].ID == taskID {
				m.tasks[i].Status = status
				m.tasks[i].UpdatedAt = time.Now()
				updated++
				break
			}
		}
	}

	fmt.Printf("MockStorage: Updated %d tasks to status %s\n", updated, status)
	return nil
}

func generateMockTasks(count int) []models.Task {
	var tasks []models.Task
	for i := 0; i < count; i++ {
		status := models.StatusPending
		if i%10 == 0 {
			status = models.StatusCompleted
		}

		task := models.Task{
			ID:      fmt.Sprintf("task-%d", i),
			Subject: fmt.Sprintf("Test Subject %d", i),
			Body: models.EmailBody{
				Text: fmt.Sprintf("Email body %d", i),
			},
			Recipients: []models.Recipient{
				{Email: fmt.Sprintf("user%d@example.com", i)},
				{Email: fmt.Sprintf("test%d@example.com", i)},
			},
			FromEmail: "noreply@myapp.com",
			Status:    status,
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
		tasks = append(tasks, task)
	}
	return tasks
}
