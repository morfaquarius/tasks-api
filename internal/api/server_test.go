package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tasks-api/internal/models"
	"tasks-api/internal/storage/memory"
)

func newTestRouter() http.Handler {
	server := NewServer(NewHandler(memory.NewMemStorage()))
	server.RegisterHandlers()
	return server.router
}

func TestCreateTaskDescriptionIsOptional(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Test task"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	newTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusCreated, response.Body)
	}

	var task models.Task
	if err := json.NewDecoder(response.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if task.Title != "Test task" || task.Description != "" {
		t.Fatalf("created task = %+v, want title set and empty description", task)
	}
	if got := response.Header().Get("Location"); got != "/tasks/1" {
		t.Errorf("Location = %q, want /tasks/1", got)
	}
}

func TestJSONBodyMustContainSingleValue(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "create",
			method: http.MethodPost,
			path:   "/tasks",
			body:   `{"title":"Test task"} {"title":"Second task"}`,
		},
		{
			name:   "update",
			method: http.MethodPut,
			path:   "/tasks/1",
			body:   `{"title":"Test task","done":false} {"title":"Second task","done":true}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			response := httptest.NewRecorder()

			newTestRouter().ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
			}
		})
	}
}

func TestInvalidTaskIDErrorIsConsistent(t *testing.T) {
	tests := []struct {
		method string
		body   string
	}{
		{method: http.MethodGet},
		{method: http.MethodPut, body: `{"title":"Test task","done":false}`},
		{method: http.MethodDelete},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/tasks/abc", strings.NewReader(test.body))
			response := httptest.NewRecorder()
			newTestRouter().ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
			var body ErrorDTO
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if body.Error != InvalidTaskID.Error() {
				t.Errorf("error = %q, want %q", body.Error, InvalidTaskID.Error())
			}
		})
	}
}

func TestTaskTitleCannotBeWhitespace(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "create",
			method: http.MethodPost,
			path:   "/tasks",
			body:   `{"title":"   "}`,
		},
		{
			name:   "update",
			method: http.MethodPut,
			path:   "/tasks/1",
			body:   `{"title":"   ","done":false}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			newTestRouter().ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
			}
		})
	}
}

func TestRoutingErrorsUseJSON(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantError  string
		wantAllow  string
	}{
		{
			name:       "unknown route",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantError:  "route not found",
		},
		{
			name:       "unsupported method",
			method:     http.MethodPatch,
			path:       "/tasks/1",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  "DELETE, GET, PUT",
		},
		{
			name:       "unsupported collection method",
			method:     http.MethodPatch,
			path:       "/tasks",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  "GET, POST",
		},
		{
			name:       "unsupported health method",
			method:     http.MethodPost,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
			wantAllow:  "GET",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()

			newTestRouter().ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if got := response.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}

			var body ErrorDTO
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if body.Error != test.wantError {
				t.Errorf("error = %q, want %q", body.Error, test.wantError)
			}
			if body.Time.IsZero() {
				t.Error("error time is empty")
			}
			if test.wantAllow != "" && response.Header().Get("Allow") != test.wantAllow {
				t.Errorf("Allow = %q, want %q", response.Header().Get("Allow"), test.wantAllow)
			}
		})
	}
}
