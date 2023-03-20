package commands

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"sync"
	"time"
)

type osCommand struct {
	cmd       *exec.Cmd
	waiters   *sync.WaitGroup
	startTime time.Time
	cancel    context.CancelFunc
}

func (o *osCommand) Pid() int {
	if o.cmd.Process != nil {
		return o.cmd.Process.Pid
	}

	return 0
}

func (o *osCommand) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	o.cancel = cancel

	if err := o.cmd.Start(); err != nil {
		return err
	}

	o.startTime = time.Now()

	o.waiters.Add(2) //nolint: gomnd

	go osSidecarSignals(ctx, o.waiters, o.cmd)

	go osSidecarGracefulShutdown(ctx, o.waiters, o.cmd)

	return nil
}

func (o *osCommand) Wait() ExecutionResult {
	err := o.cmd.Wait()
	finishTime := time.Now()

	o.cancel()
	o.waiters.Wait()

	result := ExecutionResult{
		ElapsedTime: finishTime.Sub(o.startTime),
	}

	var exitErr *exec.ExitError

	switch {
	case o.cmd.ProcessState != nil:
		result.ExitCode = o.cmd.ProcessState.ExitCode()
		result.UserTime = o.cmd.ProcessState.UserTime()
		result.SystemTime = o.cmd.ProcessState.SystemTime()
	case errors.As(err, &exitErr):
		result.ExitCode = exitErr.ExitCode()
	}

	return result
}

func New(
	command string,
	args, environ []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) Command {
	cmd := exec.Command(command, args...)

	cmd.Env = environ
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return &osCommand{
		cmd:     cmd,
		waiters: &sync.WaitGroup{},
	}
}
