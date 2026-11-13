package storage

import "tasks-api/internal/models"

type Storage interface {
	ListTasks() []models.Task
	CreateTask(task models.Task) (models.Task, error)
	GetTask(id int) (models.Task, error)
	UpdateTask(id int, task models.Task) (models.Task, error)
	DeleteTask(id int) error
}
