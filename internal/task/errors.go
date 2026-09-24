package task

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyTitle       = errors.New("empty task title")
	ErrInvalidTaskDates = errors.New("task deadline before schedule date")
	ErrTaskNotFound     = errors.New("task not found")
)

type TaskError struct {
	Op     string
	TaskID *int
	Err    error
}

func (e TaskError) Error() string {
	return fmt.Sprintf("task %d: %s: %v", *e.TaskID, e.Op, e.Err)
}

func (e TaskError) Unwrap() error {
	return e.Err
}
