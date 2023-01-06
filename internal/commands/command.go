package commands

import (
	"context"
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
}
