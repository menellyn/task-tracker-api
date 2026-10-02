package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/menellyn/task-tracker-api/internal/task"
)

type ErrorCode string

const (
	ErrorInternalServer   ErrorCode = "internal_server_error"
	ErrorInvalidID        ErrorCode = "invalid_id"
	ErrorTaskNotFound     ErrorCode = "task_not_found"
	ErrorInvalidJSON      ErrorCode = "invalid_json"
	ErrorValidation       ErrorCode = "validation_error"
	ErrorEmptyTitle       ErrorCode = "empty_title"
	ErrorInvalidTaskDates ErrorCode = "invalid_task_dates"
)

type ErrorResponse struct {
	Error   ErrorCode `json:"error"`
	Message string    `json:"message"`
}

func writeError(w http.ResponseWriter, statusCode int, err ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err,
		Message: message,
	})
}

type Service interface {
	CreateTask(taskRequest task.CreateTaskRequest) (task.TaskResponse, error)
	GetAllTasks() ([]task.TaskResponse, error)
	GetActualTasks() ([]task.TaskResponse, error)
	GetTaskByID(id int) (task.TaskResponse, error)
	DeleteTaskByID(id int) error
	MarkDoneTaskByID(id int) error
	UpdateTask(id int, taskRequest task.UpdateTaskRequest) (task.TaskResponse, error)
}

type TaskHandler struct {
	service   Service
	validator *validator.Validate
}

func NewTaskHandler(service Service, validator *validator.Validate) *TaskHandler {
	return &TaskHandler{
		service:   service,
		validator: validator,
	}
}

func (handler *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.service.GetAllTasks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(tasks)
}

func (handler *TaskHandler) GetActual(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.service.GetActualTasks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(tasks)
}

func (handler *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidID, "Invalid ID")
		return
	}
	gotTask, err := handler.service.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, ErrorTaskNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(gotTask)
}

func (handler *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var taskRequest task.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidJSON, err.Error())
		return
	}

	if err := handler.validator.Struct(taskRequest); err != nil {
		writeError(w, http.StatusUnprocessableEntity, ErrorValidation, err.Error())
		return
	}

	createdTask, err := handler.service.CreateTask(taskRequest)
	if err != nil {
		if errors.Is(err, task.ErrEmptyTitle) {
			writeError(w, http.StatusUnprocessableEntity, ErrorEmptyTitle, err.Error())

		} else if errors.Is(err, task.ErrInvalidTaskDates) {
			writeError(w, http.StatusUnprocessableEntity, ErrorInvalidTaskDates, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdTask)
}

func (handler *TaskHandler) MarkDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidID, "Invalid ID")
		return
	}
	if err := handler.service.MarkDoneTaskByID(id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, ErrorTaskNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var taskRequest task.UpdateTaskRequest

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidID, "Invalid ID")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidJSON, err.Error())
		return
	}

	if err := handler.validator.Struct(taskRequest); err != nil {
		writeError(w, http.StatusUnprocessableEntity, ErrorValidation, err.Error())
		return
	}

	updatedTask, err := handler.service.UpdateTask(id, taskRequest)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, ErrorTaskNotFound, err.Error())
		} else if errors.Is(err, task.ErrEmptyTitle) {
			writeError(w, http.StatusUnprocessableEntity, ErrorEmptyTitle, err.Error())
		} else if errors.Is(err, task.ErrInvalidTaskDates) {
			writeError(w, http.StatusUnprocessableEntity, ErrorInvalidTaskDates, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedTask)
}

func (handler *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorInvalidID, "Invalid ID")
		return
	}
	if err := handler.service.DeleteTaskByID(id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, ErrorTaskNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, ErrorInternalServer, "Internal Server Error")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
