package handler

import "github.com/menellyn/task-tracker-api/internal/task"

type FakeService struct {
	CreateTaskFunc       func(taskRequest task.CreateTaskRequest) (task.TaskResponse, error)
	UpdateTaskFunc       func(taskID int, taskRequest task.UpdateTaskRequest) (task.TaskResponse, error)
	GetAllTasksFunc      func() ([]task.TaskResponse, error)
	GetActualTasksFunc   func() ([]task.TaskResponse, error)
	GetTaskByIDFunc      func(taskID int) (task.TaskResponse, error)
	DeleteTaskByIDFunc   func(taskID int) error
	MarkDoneTaskByIDFunc func(taskID int) error
}

func (f *FakeService) CreateTask(taskRequest task.CreateTaskRequest) (task.TaskResponse, error) {
	return f.CreateTaskFunc(taskRequest)
}

func (f *FakeService) UpdateTask(id int, taskRequest task.UpdateTaskRequest) (task.TaskResponse, error) {
	return f.UpdateTaskFunc(id, taskRequest)
}

func (f *FakeService) GetAllTasks() ([]task.TaskResponse, error) {
	return f.GetAllTasksFunc()
}

func (f *FakeService) GetActualTasks() ([]task.TaskResponse, error) {
	return f.GetActualTasksFunc()
}

func (f *FakeService) GetTaskByID(id int) (task.TaskResponse, error) {
	return f.GetTaskByIDFunc(id)
}

func (f *FakeService) DeleteTaskByID(id int) error {
	return f.DeleteTaskByIDFunc(id)
}

func (f *FakeService) MarkDoneTaskByID(id int) error {
	return f.MarkDoneTaskByIDFunc(id)
}
