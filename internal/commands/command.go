package commands

import (
	"fmt"
	"time"
)

type Command interface {
	Pid() int
	Start() error
	Wait() (ExecutionResult, error)
}

type ExecutionResult struct {
	ExitCode    int
	UserTime    time.Duration
	SystemTime  time.Duration
	ElapsedTime time.Duration
}

type ExitError struct {
	Result ExecutionResult
}

func (e ExitError) Error() string {
	return fmt.Sprintf(
		"command has finished with %d in %v",
		e.Result.ExitCode,
		e.Result.ElapsedTime)
}

func (e ExitError) Code() int {
	return e.Result.ExitCode
}
