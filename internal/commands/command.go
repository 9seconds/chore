package commands

import (
	"context"
	"fmt"
	"time"
)

type Command interface {
	Pid() int
	Start(context.Context) error
	Wait() ExecutionResult
}

type ExecutionResult struct {
	ExitCode    int
	UserTime    time.Duration
	SystemTime  time.Duration
	ElapsedTime time.Duration

	err error
}

func (e ExecutionResult) Ok() bool {
	return e.Code() == 0
}

func (e ExecutionResult) Code() int {
	return e.ExitCode
}

func (e ExecutionResult) Error() string {
	return fmt.Sprintf(
		"command has finished with %d in %v",
		e.ExitCode,
		e.ElapsedTime)
}

func (e ExecutionResult) Unwrap() error {
	return e.err
}
