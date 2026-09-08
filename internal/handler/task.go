package handler

import (
	"encoding/json"
	"net/http"

	"github.com/menellyn/task-tracker-api/internal/task"
)

type TaskHandler struct {
	service *task.TaskService
}

func NewTaskHandler(service *task.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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
