package main

import (
	"fmt"
	"log"
	"os"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdRun struct {
	Namespace cli.Namespace `arg:"" help:"Script namespace."`
	Script    string        `arg:"" help:"Script name."`
	Args      []string      `arg:"" optional:"" help:"Script arguments to use."`
}

func (c *CliCmdRun) Run(ctx cli.Context) error {
	executable, err := script.New(c.Namespace.Value(), c.Script)
	if err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer executable.Cleanup()

	args, err := argparse.Parse(ctx, executable.Config.Parameters, c.Args)
	if err != nil {
		return fmt.Errorf("cannot parse arguments: %w", err)
	}

	environ := executable.Environ(ctx, args)

	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	cmd := commands.NewOS(
		ctx,
		executable,
		environ,
		args.Positional,
		os.Stdin,
		os.Stdout,
		os.Stderr)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start command: %w", err)
	}

	log.Printf("command %s has started as %d", executable, cmd.Pid())

	result, err := cmd.Wait()
	if err != nil {
		return fmt.Errorf("cannot correctly finish command: %w", err)
	}

	log.Printf("command %d exit with exit code %d", cmd.Pid(), result.ExitCode)
	log.Printf(
		"command %d times: user=%v, sys=%v, real=%v",
		cmd.Pid(),
		result.UserTime,
		result.SystemTime,
		result.ElapsedTime)

	return commands.ExitError{
		Result: result,
	}
}
