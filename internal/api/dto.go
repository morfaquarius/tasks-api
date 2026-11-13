package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}
	if t.Description == "" {
		return errors.New("description is empty")
	}
	return nil
}

type UpdateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
}

func (t UpdateTaskDTO) ValidateForUpdate() error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	if t.Description == "" {
		return errors.New("description cannot be empty")
	}
	if t.Done == nil {
		return errors.New("field 'done' is required")
	}
	return nil
}

type ErrorDTO struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

func WriteJSONError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errDTO := ErrorDTO{
		Error: err.Error(),
		Time:  time.Now(),
	}

	json.NewEncoder(w).Encode(errDTO)
}

var InvalidTaskID = errors.New("invalid task id")
