package errors

import "errors"

var (
	ErrProfilingNotFound    = errors.New("profiling not found")
	ErrProfilingFilledToday = errors.New("today's profiling has been filled")
)
