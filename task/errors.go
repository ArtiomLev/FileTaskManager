package task

import "errors"

var (
	ErrTaskExists             = errors.New("task already exists")
	ErrTaskNotExists          = errors.New("task does not exist")
	ErrActiveTaskNotExists    = errors.New("active task does not exist")
	ErrCompletedTaskNotExists = errors.New("completed task does not exist")
	ErrInvalidContract        = errors.New("invalid or empty contract number")
	ErrInvalidName            = errors.New("invalid or empty task name")
	ErrNoFiles                = errors.New("no files provided")
	ErrCannotMoveTempToTask   = errors.New("cannot move temp files to task directory")
	ErrTaskCompleted          = errors.New("task already completed")
	ErrTaskActive             = errors.New("task is active")
)
