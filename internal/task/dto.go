package task

import (
	"time"
)

type CreateTaskRequest struct {
	Title        string     `json:"title" validate:"required,min=1,max=250"`
	Description  *string    `json:"description" validate:"omitempty,max=2000"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
}

type UpdateTaskRequest struct {
	Title        *string    `json:"title" validate:"omitempty,min=1,max=250"`
	Description  *string    `json:"description" validate:"omitempty,max=2000"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         *bool      `json:"done"`
}

type TaskResponse struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         bool       `json:"done"`
}

func ToTaskResponse(t Task) TaskResponse {
	return TaskResponse{
		ID:           t.ID,
		Title:        t.Title,
		Description:  &t.Description,
		ScheduleDate: t.ScheduleDate,
		Deadline:     t.Deadline,
		Done:         t.Done,
	}
}

func ToTaskModel()
