package task

import (
	"encoding/json"
	"time"
)

type CreateTaskRequest struct {
	Title        string     `json:"title" validate:"required,min=1,max=250"`
	Description  *string    `json:"description" validate:"omitempty,max=2000"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
}

func (r *CreateTaskRequest) ToTaskModel() Task {
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

	scheduleDateSet bool
	deadlineSet     bool
}

func (r *UpdateTaskRequest) UnmarshalJSON(data []byte) error {
	type Alias UpdateTaskRequest

	var aux Alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*r = UpdateTaskRequest(aux)

	var fields map[string]json.RawMessage

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	_, r.scheduleDateSet = fields["scheduleDate"]
	_, r.deadlineSet = fields["deadline"]

	return nil
}

func (r *UpdateTaskRequest) ToMap() map[string]interface{} {
	updateData := make(map[string]interface{})
	if r.Title != nil {
		updateData["title"] = *r.Title
	}
	if r.Description != nil {
		updateData["description"] = r.Description
	}
	if r.scheduleDateSet {
		updateData["schedule_date"] = r.ScheduleDate
	}
	if r.deadlineSet {
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
