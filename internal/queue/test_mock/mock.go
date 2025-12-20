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

func NewMockStorage(taskCount, emailsPerTask int) *MockStorage {
	return &MockStorage{
		tasks: generateMockTasks(taskCount, emailsPerTask),
	}
}

func (m *MockStorage) Enqueue(ctx context.Context, task models.Task) error {
	return nil
}

func (m *MockStorage) Dequeue(ctx context.Context) (models.Task, error) {
	return m.tasks[0], nil
}

func (m *MockStorage) DequeueBatch(ctx context.Context, batchSize int) ([]models.Task, error) {
	return m.tasks, nil
}

func (m *MockStorage) CheckStatus(ctx context.Context) error {
	return nil
}

func (m *MockStorage) Close() {

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

func generateMockTasks(taskCount, emailsPerTask int) []models.Task {
	var tasks []models.Task

	status := models.StatusPending
	for taskIndex := 0; taskIndex < taskCount; taskIndex++ {
		recipients := make([]models.Recipient, 0, emailsPerTask)
		for emailIndex := 0; emailIndex < emailsPerTask; emailIndex++ {
			recipient := models.Recipient{
				Email: fmt.Sprintf("user%d.task%d@example.com", emailIndex, taskIndex),
			}
			recipients = append(recipients, recipient)
		}

		task := models.Task{
			ID:      fmt.Sprintf("task-%03d", taskIndex),
			Subject: fmt.Sprintf("Task #%d - %s Status", taskIndex, status),
			Body: models.EmailBody{
				Text: fmt.Sprintf("This is email body for task %d\nTotal recipients: %d\nGenerated at: %s",
					taskIndex, emailsPerTask, time.Now().Format(time.RFC3339)),
				HTML: fmt.Sprintf("<html><body><h1>Task %d</h1><p>Total recipients: %d</p></body></html>",
					taskIndex, emailsPerTask),
			},
			Recipients: recipients,
			FromEmail:  "noreply@myapp.com",
			Status:     status,
			CreatedAt:  time.Now().Add(-time.Duration(taskCount-taskIndex) * time.Hour),
			UpdatedAt:  time.Now().Add(-time.Duration(taskCount-taskIndex) * time.Hour),
		}
		tasks = append(tasks, task)
	}

	// for i := range taskCount {
	// 	status := models.StatusPending

	// 	task := models.Task{
	// 		ID:      fmt.Sprintf("task-%d", i),
	// 		Subject: fmt.Sprintf("Test Subject %d", i),
	// 		Body: models.EmailBody{
	// 			Text: fmt.Sprintf("Email body %d", i),
	// 		},
	// 		Recipients: []models.Recipient{
	// 			{Email: fmt.Sprintf("user%d@example.com", i)},
	// 			// {Email: fmt.Sprintf("test%d@example.com", i)},
	// 		},
	// 		FromEmail: "noreply@myapp.com",
	// 		Status:    status,
	// 		CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
	// 		UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
	// 	}
	// 	tasks = append(tasks, task)
	// }

	return tasks
}
