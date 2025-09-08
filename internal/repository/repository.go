package repository

import "Gonoty/internal/models"

type TaskRepository interface {
	AddTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
	LoadTasks() ([]models.Task, error)
}
