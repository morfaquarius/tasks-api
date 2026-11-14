package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t TaskDTO) ValidateForCreate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is empty")
	}
	return nil
}

type UpdateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
}

func (t UpdateTaskDTO) ValidateForUpdate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title cannot be empty")
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
	writeJSONResponse(w, nil, statusCode, ErrorDTO{
		Error: err.Error(),
		Time:  time.Now(),
	})
}

func writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, response any) {
	body, err := json.Marshal(response)
	if err != nil {
		method, path := "", ""
		if r != nil {
			method, path = r.Method, r.URL.Path
		}
		log.Printf("error encoding %s %s response: %v", method, path, err)
		if statusCode != http.StatusInternalServerError {
			WriteJSONError(w, errors.New("internal server error"), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(append(body, '\n')); err != nil {
		method, path := "", ""
		if r != nil {
			method, path = r.Method, r.URL.Path
		}
		log.Printf("error writing %s %s response: %v", method, path, err)
	}
}

var InvalidTaskID = errors.New("invalid task id")
