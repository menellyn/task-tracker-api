package task

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyTitle = errors.New("empty task title")

type Repository interface {
	Add(title string) (Task, error)
	GetAll() ([]Task, error)
	GetActual() ([]Task, error)
	GetByID(id int) (Task, error)
	MarkDone(id int) error
	Update(task Task) (Task, error)
	Delete(id int) error
}
type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	return s.repo.Add(title)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetActualTasks() ([]Task, error) {
	return s.repo.GetActual()
}

func (s *TaskService) GetTaskByID(id int) (Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return Task{}, fmt.Errorf("get task by id %d: %w", id, err)
	}
	return task, nil
}

func (s *TaskService) DeleteTaskByID(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("delete task by id %d: %w", id, err)
	}
	return nil
}

func (s *TaskService) MarkDoneTaskByID(id int) error {
	err := s.repo.MarkDone(id)
	if err != nil {
		return fmt.Errorf("mark done task %d: %w", id, err)
	}
	return nil
}

func (s *TaskService) Update(task Task) (Task, error) {
	updatedTask, err := s.repo.Update(task)
	if err != nil {
		return Task{}, fmt.Errorf("update task %d: %w", task.ID, err)
	}
	return updatedTask, nil
}
