package task

import "time"

type Task struct {
	ID           int
	Title        string
	Description  *string
	ScheduleDate *time.Time
	Deadline     *time.Time
	Done         bool
}
