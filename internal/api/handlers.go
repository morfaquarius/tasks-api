package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"

	"github.com/gorilla/mux"
)

type Handlers struct{ store storage.Storage }

func NewHandler(s storage.Storage) *Handlers { return &Handlers{store: s} }

func (h *Handlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		WriteJSONError(w, err, http.StatusBadRequest)
		return
	}
	if err := taskDTO.ValidateForCreate(); err != nil {
		WriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	newTask := models.NewTask(taskDTO.Title, taskDTO.Description)
	createdTask, err := h.store.CreateTask(newTask)
	if err != nil {
		WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		log.Printf("error encoding create response: %v", err)
		return
	}
}

func (h *Handlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	reqId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqId)
	if err != nil {
		WriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	getTask, err := h.store.GetTask(id)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			WriteJSONError(w, err, http.StatusNotFound)
		} else {
			WriteJSONError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getTask); err != nil {
		log.Printf("error encoding create response: %v", err)
		return
	}
}

func (h *Handlers) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.ListTasks()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("error encoding create response: %v", err)
		return
	}
}

func (h *Handlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	reqId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqId)
	if err != nil {
		WriteJSONError(w, InvalidTaskID, http.StatusBadRequest)
		return
	}

	var taskDTO UpdateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		WriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForUpdate(); err != nil {
		WriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	taskToUpdate := models.Task{
		ID:          id,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		Done:        *taskDTO.Done,
	}

	updatedTask, err := h.store.UpdateTask(id, taskToUpdate)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			WriteJSONError(w, err, http.StatusNotFound)
		} else {
			WriteJSONError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		log.Printf("error encoding create response: %v", err)
		return
	}
}

func (h *Handlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	reqId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqId)
	if err != nil {
		WriteJSONError(w, InvalidTaskID, http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteTask(id); err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			WriteJSONError(w, err, http.StatusNotFound)
		} else {
			WriteJSONError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
