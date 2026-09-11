package task

import "time"

type CreateTaskRequest struct {
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
}

type UpdateTaskRequest struct {
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         *bool      `json:"done"`
}

type TaskResponse struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         bool       `json:"done"`
}
