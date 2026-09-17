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

func (r CreateTaskRequest) ToTaskModel() Task {
	return Task{
		Title:        r.Title,
		Description:  r.Description,
		ScheduleDate: r.ScheduleDate,
		Deadline:     r.Deadline,
	}
}

type UpdateTaskRequest struct {
	Title        *string    `json:"title" validate:"omitempty,min=1,max=250"`
	Description  *string    `json:"description" validate:"omitempty,max=2000"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         *bool      `json:"done"`
}

func (r UpdateTaskRequest) ToMap() map[string]interface{} {
	updateData := make(map[string]interface{})
	if r.Title != nil {
		updateData["title"] = *r.Title
	}
	if r.Description != nil {
		updateData["description"] = r.Description
	}
	if r.ScheduleDate != nil {
		updateData["schedule_date"] = r.ScheduleDate
	}
	if r.Deadline != nil {
		updateData["deadline"] = r.Deadline
	}
	if r.Done != nil {
		updateData["done"] = *r.Done
	}
	return updateData
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
		Description:  t.Description,
		ScheduleDate: t.ScheduleDate,
		Deadline:     t.Deadline,
		Done:         t.Done,
	}
}
