package repository

import (
	"Gonoty/internal/models"
	"errors"
	"fmt"
	"os"
	"sync"
)

type FileRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewFileRepository(filePath string) (*FileRepository, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.New("failed to create file repository")
	}
	fileRepository := &FileRepository{file: file}
	return fileRepository, nil
}

func (r *FileRepository) AddTask(task *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Fprintf(r.file, "%v\n", task)
	return nil
}

func (r *FileRepository) LoadTasks() ([]models.Task, error) {
	var tasks []models.Task
	return nil, nil
}

func (r *FileRepository) UpdateTask(task *models.Task) error {
	return nil
}

func (r *FileRepository) DeleteTask(task *models.Task) error {
	return nil
}

func (r *FileRepository) Close() {
	r.file.Close()
}
