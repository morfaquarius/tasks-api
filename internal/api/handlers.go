package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"

	"github.com/gorilla/mux"
)

type Handlers struct{ store storage.Storage }

func NewHandler(s storage.Storage) *Handlers { return &Handlers{store: s} }

func decodeJSONBody(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("request body must contain a single JSON value")
		}
		return err
	}
	return nil
}

func (h *Handlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := decodeJSONBody(r, &taskDTO); err != nil {
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

	w.Header().Set("Location", fmt.Sprintf("/tasks/%d", createdTask.ID))
	writeJSONResponse(w, r, http.StatusCreated, createdTask)
}

func (h *Handlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	reqId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqId)
	if err != nil {
		WriteJSONError(w, InvalidTaskID, http.StatusBadRequest)
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

	writeJSONResponse(w, r, http.StatusOK, getTask)
}

func (h *Handlers) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.ListTasks()

	writeJSONResponse(w, r, http.StatusOK, tasks)
}

func (h *Handlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	reqId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(reqId)
	if err != nil {
		WriteJSONError(w, InvalidTaskID, http.StatusBadRequest)
		return
	}

	var taskDTO UpdateTaskDTO
	if err := decodeJSONBody(r, &taskDTO); err != nil {
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

	writeJSONResponse(w, r, http.StatusOK, updatedTask)
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

func (h *Handlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, r, http.StatusOK, map[string]string{"status": "ok"})
}
