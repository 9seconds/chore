package commands

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"

	"github.com/9seconds/chore/internal/script"
)

var ErrNotStarted = errors.New("process not started")

type osCommand struct {
	cmd       *exec.Cmd
	ctx       context.Context
	ctxCancel context.CancelFunc

	finishErr error

	startTime  time.Time
	finishTime time.Time
	startOnce  sync.Once
	finishOnce sync.Once

	sigChan chan os.Signal
	sigDone chan struct{}
}

func (o *osCommand) Pid() int {
	if o.cmd.Process != nil {
		return o.cmd.Process.Pid
	}

	return 0
}

func (o *osCommand) Start() error {
	var (
		err  error
		boot bool
	)

	o.startOnce.Do(func() {
		err = o.cmd.Start()
		o.startTime = time.Now()
		boot = true

		signal.Notify(o.sigChan, SignalsToRelay...)

		go func() {
			defer func() {
				signal.Stop(o.sigChan)
				close(o.sigDone)
			}()

			for {
				select {
				case <-o.ctx.Done():
					return
				case sig := <-o.sigChan:
					switch {
					case o.cmd.Process == nil:
						continue
					case o.cmd.ProcessState != nil && o.cmd.ProcessState.Exited():
						return
					}

					log.Printf("!end %s to %d", sig, o.cmd.Process.Pid)

					if err := o.cmd.Process.Signal(sig); err != nil {
						log.Printf("cannot send %v to process %d: %v", sig, o.Pid(), err)
					}
				}
			}
		}()
	})

	if boot {
		return err
	}

	return o.cmd.Start()
}

func (o *osCommand) Wait() (ExecutionResult, error) {
	result := ExecutionResult{}

	if o.startTime.IsZero() {
		return result, ErrNotStarted
	}

	o.finishOnce.Do(func() {
		o.finishErr = o.cmd.Wait()
		o.finishTime = time.Now()
		o.ctxCancel()
		<-o.sigDone
	})

	if o.cmd.ProcessState != nil {
		result.ExitCode = o.cmd.ProcessState.ExitCode()
		result.UserTime = o.cmd.ProcessState.UserTime()
		result.SystemTime = o.cmd.ProcessState.SystemTime()
		result.ElapsedTime = o.finishTime.Sub(o.startTime)
	}

	var exitErr *exec.ExitError

	if errors.As(o.finishErr, &exitErr) {
		result.ExitCode = exitErr.ExitCode()

		return result, nil
	}

	return result, o.finishErr
}

func NewOS(ctx context.Context, script *script.Script,
	environ, args []string,
	stdin io.Reader, stdout, stderr io.Writer,
) Command {
	ctx, cancel := context.WithCancel(ctx)

	cmdLine := []string{script.Path()}
	cmdLine = append(cmdLine, args...)

	cmd := exec.CommandContext(ctx, cmdLine[0], cmdLine[1:]...)

	cmd.Env = append(os.Environ(), environ...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return &osCommand{
		cmd:       cmd,
		ctx:       ctx,
		ctxCancel: cancel,
		sigChan:   make(chan os.Signal, 1),
		sigDone:   make(chan struct{}),
	}
}
