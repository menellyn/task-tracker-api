package task

import "time"

type Task struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ScheduleDate *time.Time `json:"scheduleDate"`
	Deadline     *time.Time `json:"deadline"`
	Done         bool       `json:"done"`
}
