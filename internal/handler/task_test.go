package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/menellyn/task-tracker-api/internal/task"
)

func TestTaskHandler_CreateTask(t *testing.T) {
	tests := []struct {
		name           string
		initBody       string
		createTaskFunc func(request task.CreateTaskRequest) (task.TaskResponse, error)
		wantResponse   task.TaskResponse
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:     "success",
			initBody: `{"title":"Buy milk"}`,
			createTaskFunc: func(request task.CreateTaskRequest) (task.TaskResponse, error) {
				return task.TaskResponse{
						ID:    1,
						Title: request.Title,
						Done:  false,
					},
					nil
			},
			wantResponse: task.TaskResponse{
				ID:    1,
				Title: "Buy milk",
				Done:  false,
			},
			wantErr:        false,
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "empty title",
			initBody:       `{"title":""}`,
			wantErr:        true,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "invalid dates",
			initBody: `{"title":"Buy milk", "scheduleDate":"2026-09-25T00:00:00Z", "deadline":"2026-09-20T00:00:00Z"}`,
			createTaskFunc: func(request task.CreateTaskRequest) (task.TaskResponse, error) {
				return task.TaskResponse{}, task.ErrInvalidTaskDates
			},
			wantErr:        true,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid json",
			initBody:       `{"title":"Buy milk"`,
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				CreateTaskFunc: tt.createTaskFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)

			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader(tt.initBody),
			)
			rr := httptest.NewRecorder()

			handler.Create(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
			}

			if !tt.wantErr {
				gotResponse := task.TaskResponse{}
				err := json.NewDecoder(rr.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("unexpected response body: got %+v want %+v", gotResponse, tt.wantResponse)
				}
			}
		})
	}

}

func TestTaskHandler_UpdateTask(t *testing.T) {
	tests := []struct {
		name           string
		targetPath     string
		pathID         string
		updateBody     string
		updateTaskFunc func(id int, request task.UpdateTaskRequest) (task.TaskResponse, error)
		wantResponse   task.TaskResponse
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:       "success",
			targetPath: "/tasks/1",
			pathID:     "1",
			updateBody: `{"title":"Buy milk!"}`,
			updateTaskFunc: func(id int, request task.UpdateTaskRequest) (task.TaskResponse, error) {
				title := request.Title
				return task.TaskResponse{
					ID:    id,
					Title: *title,
					Done:  false,
				}, nil
			},
			wantResponse: task.TaskResponse{
				ID:    1,
				Title: "Buy milk!",
				Done:  false,
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name:       "not found",
			targetPath: "/tasks/1",
			pathID:     "1",
			updateBody: `{"title":"Buy milk!"}`,
			updateTaskFunc: func(id int, request task.UpdateTaskRequest) (task.TaskResponse, error) {
				return task.TaskResponse{}, task.ErrTaskNotFound
			},
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "invalid json",
			targetPath:     "/tasks/1",
			pathID:         "1",
			updateBody:     `{"title":"Buy milk"`,
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:       "invalid dates",
			targetPath: "/tasks/1",
			pathID:     "1",
			updateBody: `{"title":"Buy milk", "scheduleDate":"2026-09-25T00:00:00Z", "deadline":"2026-09-20T00:00:00Z"}`,
			updateTaskFunc: func(id int, request task.UpdateTaskRequest) (task.TaskResponse, error) {
				return task.TaskResponse{}, task.ErrInvalidTaskDates
			},
			wantErr:        true,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty title",
			targetPath:     "/tasks/1",
			pathID:         "1",
			updateBody:     `{"title":""}`,
			wantErr:        true,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid id",
			targetPath:     "/tasks/one",
			pathID:         "one",
			updateBody:     `{"title":"Buy milk"}`,
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				UpdateTaskFunc: tt.updateTaskFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)

			req := httptest.NewRequest(http.MethodPatch, tt.targetPath, strings.NewReader(tt.updateBody))
			req.SetPathValue("id", tt.pathID)
			rr := httptest.NewRecorder()

			handler.Update(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
			}

			if !tt.wantErr {
				gotResponse := task.TaskResponse{}
				err := json.NewDecoder(rr.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("unexpected response body: got %+v want %+v", gotResponse, tt.wantResponse)
				}
			}

		})
	}
}

func TestTaskHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		getAllTaskFunc func() ([]task.TaskResponse, error)
		wantResponse   []task.TaskResponse
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "success",
			getAllTaskFunc: func() ([]task.TaskResponse, error) {
				return []task.TaskResponse{
					{
						ID:    1,
						Title: "Buy milk",
						Done:  false,
					},
					{
						ID:    2,
						Title: "Wash dishes",
						Done:  true,
					},
				}, nil
			},
			wantResponse: []task.TaskResponse{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
				{
					ID:    2,
					Title: "Wash dishes",
					Done:  true,
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "database error",
			getAllTaskFunc: func() ([]task.TaskResponse, error) {
				return nil, errors.New("database error")
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				GetAllTasksFunc: tt.getAllTaskFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)

			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			rr := httptest.NewRecorder()

			handler.GetAll(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}
			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
			}

			if !tt.wantErr {
				gotResponse := []task.TaskResponse{}
				err := json.NewDecoder(rr.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("unexpected response body: got %+v want %+v", gotResponse, tt.wantResponse)
				}
			}
		})
	}
}

func TestTaskHandler_GetActual(t *testing.T) {
	tests := []struct {
		name               string
		getActualTasksFunc func() ([]task.TaskResponse, error)
		wantResponse       []task.TaskResponse
		wantErr            bool
		wantStatusCode     int
	}{
		{
			name: "success",
			getActualTasksFunc: func() ([]task.TaskResponse, error) {
				return []task.TaskResponse{
					{
						ID:    1,
						Title: "Buy milk",
						Done:  false,
					},
				}, nil
			},
			wantResponse: []task.TaskResponse{
				{
					ID:    1,
					Title: "Buy milk",
					Done:  false,
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "database error",
			getActualTasksFunc: func() ([]task.TaskResponse, error) {
				return nil, errors.New("database error")
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				GetActualTasksFunc: tt.getActualTasksFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)

			req := httptest.NewRequest(http.MethodGet, "/tasks/actual", nil)
			rr := httptest.NewRecorder()

			handler.GetActual(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}
			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
			}
			if !tt.wantErr {
				gotResponse := []task.TaskResponse{}
				err := json.NewDecoder(rr.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("unexpected response body: got %+v want %+v", gotResponse, tt.wantResponse)
				}
			}
		})
	}
}

func TestTaskHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		idPath         string
		getByIDFunc    func(taskID int) (task.TaskResponse, error)
		wantResponse   task.TaskResponse
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:   "success",
			idPath: "1",
			getByIDFunc: func(taskID int) (task.TaskResponse, error) {
				return task.TaskResponse{
					ID:    taskID,
					Title: "Buy milk",
					Done:  false,
				}, nil
			},
			wantResponse: task.TaskResponse{
				ID:    1,
				Title: "Buy milk",
				Done:  false,
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "not found",
			idPath: "1",
			getByIDFunc: func(taskID int) (task.TaskResponse, error) {
				return task.TaskResponse{}, task.ErrTaskNotFound
			},
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:   "database error",
			idPath: "1",
			getByIDFunc: func(taskID int) (task.TaskResponse, error) {
				return task.TaskResponse{}, errors.New("database error")
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "invalid id",
			idPath:         "one",
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				GetTaskByIDFunc: tt.getByIDFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)

			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.idPath, nil)
			req.SetPathValue("id", tt.idPath)
			rr := httptest.NewRecorder()
			handler.GetByID(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}
			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
			}

			if !tt.wantErr {
				gotResponse := task.TaskResponse{}
				err := json.NewDecoder(rr.Body).Decode(&gotResponse)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("unexpected response body: got %+v want %+v", gotResponse, tt.wantResponse)
				}
			}
		})
	}
}

func TestTaskHandler_MarkDone(t *testing.T) {
	tests := []struct {
		name             string
		idPath           string
		markDoneTaskFunc func(taskID int) error
		wantErr          bool
		wantStatusCode   int
	}{
		{
			name:   "success",
			idPath: "1",
			markDoneTaskFunc: func(taskID int) error {
				return nil
			},
			wantErr:        false,
			wantStatusCode: http.StatusNoContent,
		},
		{
			name:   "database error",
			idPath: "1",
			markDoneTaskFunc: func(taskID int) error {
				return errors.New("database error")
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "invalid id",
			idPath:         "one",
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				MarkDoneTaskByIDFunc: tt.markDoneTaskFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)
			req := httptest.NewRequest(http.MethodPatch, "/tasks/"+tt.idPath+"/done", nil)
			req.SetPathValue("id", tt.idPath)
			rr := httptest.NewRecorder()
			handler.MarkDone(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}
			if tt.wantErr {
				if rr.Header().Get("Content-Type") != "application/json" {
					t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
				}
			}
		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		idPath         string
		deleteTaskFunc func(taskID int) error
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:   "success",
			idPath: "1",
			deleteTaskFunc: func(taskID int) error {
				return nil
			},
			wantErr:        false,
			wantStatusCode: http.StatusNoContent,
		},
		{
			name:   "database error",
			idPath: "1",
			deleteTaskFunc: func(taskID int) error {
				return errors.New("database error")
			},
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "invalid id",
			idPath:         "one",
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FakeService{
				DeleteTaskByIDFunc: tt.deleteTaskFunc,
			}
			validate := validator.New()
			handler := NewTaskHandler(service, validate)
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.idPath, nil)
			req.SetPathValue("id", tt.idPath)
			rr := httptest.NewRecorder()
			handler.Delete(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatusCode)
			}
			if tt.wantErr {
				if rr.Header().Get("Content-Type") != "application/json" {
					t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
				}
			}
		})
	}
}
