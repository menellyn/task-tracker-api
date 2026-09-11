package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/menellyn/task-tracker-api/internal/task"
)

type TaskHandler struct {
	service   *task.TaskService
	validator *validator.Validate
}

func NewTaskHandler(service *task.TaskService, validator *validator.Validate) *TaskHandler {
	return &TaskHandler{
		service:   service,
		validator: validator,
	}
}

func (handler *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.service.GetAllTasks()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *TaskHandler) GetActual(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.service.GetActualTasks()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	gotTask, err := handler.service.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gotTask); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var taskToCreate task.Task
	if err := json.NewDecoder(r.Body).Decode(&taskToCreate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdTask, err := handler.service.CreateTask(taskToCreate.Title)
	if err != nil {
		if errors.Is(err, task.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)

		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *TaskHandler) MarkDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := handler.service.MarkDoneTaskByID(id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var taskToUpdate task.Task
	if err := json.NewDecoder(r.Body).Decode(&taskToUpdate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedTask, err := handler.service.Update(taskToUpdate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (handler *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := handler.service.DeleteTaskByID(id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
