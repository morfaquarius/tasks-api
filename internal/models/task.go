package models

import (
	"errors"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func NewTask(title string, description string) Task {
	return Task{
		Title:       title,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

var ErrTaskNotFound = errors.New("task not found")
