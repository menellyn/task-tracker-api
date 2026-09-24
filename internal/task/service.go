package task

import (
	"strings"
)

type Repository interface {
	Add(task Task) (Task, error)
	GetAll() ([]Task, error)
	GetActual() ([]Task, error)
	GetByID(id int) (Task, error)
	MarkDone(id int) error
	Update(id int, updateData map[string]interface{}) (Task, error)
	Delete(id int) error
}
type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(taskRequest CreateTaskRequest) (TaskResponse, error) {
	taskRequest.Title = strings.TrimSpace(taskRequest.Title)
	if taskRequest.Title == "" {
		return TaskResponse{}, ErrEmptyTitle
	}
	if taskRequest.Deadline != nil &&
		taskRequest.ScheduleDate != nil &&
		!taskRequest.Deadline.After(*taskRequest.ScheduleDate) {
		return TaskResponse{}, ErrInvalidTaskDates
	}

	task, err := s.repo.Add(taskRequest.ToTaskModel())
	if err != nil {
		return TaskResponse{}, err
	}
	return ToTaskResponse(task), nil
}

func (s *TaskService) GetAllTasks() ([]TaskResponse, error) {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return []TaskResponse{}, err
	}

	tasksResponse := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		tasksResponse[i] = ToTaskResponse(task)
	}

	return tasksResponse, nil
}

func (s *TaskService) GetActualTasks() ([]TaskResponse, error) {
	tasks, err := s.repo.GetActual()
	if err != nil {
		return []TaskResponse{}, err
	}

	tasksResponse := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		tasksResponse[i] = ToTaskResponse(task)
	}

	return tasksResponse, nil
}

func (s *TaskService) GetTaskByID(id int) (TaskResponse, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return TaskResponse{},
			TaskError{
				Op:     "get task",
				TaskID: &id,
				Err:    err,
			}
	}
	return ToTaskResponse(task), nil
}

func (s *TaskService) DeleteTaskByID(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return TaskError{
			Op:     "delete task",
			TaskID: &id,
			Err:    err,
		}
	}
	return nil
}

func (s *TaskService) MarkDoneTaskByID(id int) error {
	err := s.repo.MarkDone(id)
	if err != nil {
		return TaskError{
			Op:     "mark done task",
			TaskID: &id,
			Err:    err,
		}
	}
	return nil
}

func (s *TaskService) UpdateTask(id int, taskRequest UpdateTaskRequest) (TaskResponse, error) {
	if taskRequest.Title != nil {
		if strings.TrimSpace(*taskRequest.Title) == "" {
			return TaskResponse{}, ErrEmptyTitle
		}
	}
	if taskRequest.Deadline != nil &&
		taskRequest.ScheduleDate != nil &&
		!taskRequest.Deadline.After(*taskRequest.ScheduleDate) {
		return TaskResponse{}, ErrInvalidTaskDates
	}
	if taskRequest.Deadline != nil || taskRequest.ScheduleDate != nil {
		task, err := s.repo.GetByID(id)
		if err != nil {
			return TaskResponse{},
				TaskError{
					Op:     "update task",
					TaskID: &id,
					Err:    err,
				}
		}

		scheduleDate := task.ScheduleDate
		deadline := task.Deadline

		if taskRequest.Deadline != nil {
			deadline = taskRequest.Deadline
		}
		if taskRequest.ScheduleDate != nil {
			scheduleDate = taskRequest.ScheduleDate
		}

		if deadline != nil &&
			scheduleDate != nil &&
			!deadline.After(*scheduleDate) {
			return TaskResponse{}, ErrInvalidTaskDates
		}
	}

	updatedTask, err := s.repo.Update(id, taskRequest.ToMap())
	if err != nil {
		return TaskResponse{}, TaskError{
			Op:     "update task",
			TaskID: &id,
			Err:    err,
		}
	}

	return ToTaskResponse(updatedTask), nil
}
