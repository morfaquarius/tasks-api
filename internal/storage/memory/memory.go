package memory

import (
	"sync"
	"tasks-api/internal/models"
)

type MemStorage struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *MemStorage) CreateTask(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID

	s.tasks[task.ID] = task
	s.nextID++

	return task, nil
}

func (s *MemStorage) ListTasks() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmp := make([]models.Task, 0, len(s.tasks))

	for _, v := range s.tasks {
		tmp = append(tmp, v)
	}

	return tmp
}

func (s *MemStorage) GetTask(id int) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, models.ErrTaskNotFound
	}

	return task, nil
}

func (s *MemStorage) UpdateTask(id int, task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskOld, ok := s.tasks[id]
	if !ok {
		return models.Task{}, models.ErrTaskNotFound
	}

	taskOld.Title = task.Title
	taskOld.Description = task.Description
	taskOld.Done = task.Done

	s.tasks[id] = taskOld

	return taskOld, nil
}

func (s *MemStorage) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	if !ok {
		return models.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
