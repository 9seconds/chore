package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/9seconds/chore/chorelib/env"
	"github.com/9seconds/chore/chorelib/script"
)

type scriptExitError struct {
	code int
}

func (s scriptExitError) Error() string {
	return fmt.Sprintf("exit with code %d", s.code)
}

type CliCmdRun struct {
	Namespace CliNamespace `arg:"" help:"Script namespace."`
	Script    string       `arg:"" help:"Script name."`
	Args      []string     `arg:"" optional:"" help:"Script arguments to use."`
}

func (c *CliCmdRun) Run(ctx Context) error {
	executable := script.Script{
		Namespace:  c.Namespace.Value,
		Executable: c.Script,
	}

	if err := executable.IsValid(); err != nil {
		return fmt.Errorf("script is invalid: %w", err)
	}

	if err := executable.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer executable.Cleanup()

	if err := os.MkdirAll(executable.PersistentDir(), 0750); err != nil {
		return fmt.Errorf(
			"cannot ensure persistent dir %s: %w",
			executable.PersistentDir(),
			err)
	}

	environ, err := env.MakeEnviron(ctx, executable, c.Args)
	if err != nil {
		return fmt.Errorf("cannot generate environment variables: %w", err)
	}

	log.Printf(
		"run: namespace=%s, executable=%s, args=%v",
		executable.Namespace,
		executable.Executable,
		c.Args)
	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	cmd := exec.CommandContext(ctx, executable.Path())
	cmd.Env = append(os.Environ(), environ...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	timeStart := time.Now()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command cannot start: %w", err)
	}

	log.Printf("process %d has started", cmd.Process.Pid)

	relaySignals(ctx, cmd)

	err = cmd.Wait()
	timeStop := time.Now()

	exitCode := 0

	if err != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	log.Printf("sys time %v", cmd.ProcessState.SystemTime())
	log.Printf("user time %v", cmd.ProcessState.UserTime())
	log.Printf("elapsed time %v", timeStop.Sub(timeStart))
	log.Printf("process exited with %d", exitCode)

	return scriptExitError{
		code: exitCode,
	}
}
