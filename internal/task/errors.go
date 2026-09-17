package task

import "errors"

var (
	ErrEmptyTitle       = errors.New("empty task title")
	ErrInvalidTaskDates = errors.New("task deadline before schedule date")
	ErrTaskNotFound     = errors.New("task not found")
)
