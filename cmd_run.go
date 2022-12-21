package main

import (
	"fmt"
	"log"
	"os"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/filelock"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdRun struct {
	Timeout cli.Timeout `short:"t" help:"Execute with a given timeout. Number mean seconds. Also can pass duration. Default is no timeout."`
	Lock    cli.Lock    `short:"l" help:"A path to a lock to acquire before execution. Prefix 's:' means shared lock, 'x:' - exclusive (default). '.' means a path to the script itself. Default is no lock."`

	Namespace cli.Namespace `arg:"" help:"Prefix of the script namespace."`
	Script    string        `arg:"" help:"Prefix of the script name."`
	Args      []string      `arg:"" optional:"" help:"Script arguments to use."`
}

func (c *CliCmdRun) Run(ctx cli.Context) error {
	scr, err := script.FindScript(c.Namespace.Value(), c.Script)
	if err != nil {
		return fmt.Errorf("cannot find out script: %w", err)
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	if c.Timeout.Value() != 0 {
		ctx = ctx.WithTimeout(c.Timeout.Value())
	}

	lock, err := filelock.New(c.Lock.LockMode(), c.Lock.Path(scr.Path()))
	if err != nil {
		return fmt.Errorf("cannot initialize lock: %w", err)
	}

	ctx = ctx.WithLock(lock)

	if err := ctx.Start(); err != nil {
		return fmt.Errorf("cannot start context: %w", err)
	}

	conf := scr.Config()

	args, err := argparse.Parse(ctx, c.Args, conf.Flags, conf.Parameters)
	if err != nil {
		return fmt.Errorf("cannot parse arguments: %w", err)
	}

	environ := scr.Environ(ctx, args)

	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	cmd := commands.NewOS(
		ctx,
		scr,
		environ,
		args.Positional,
		os.Stdin,
		os.Stdout,
		os.Stderr)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start command: %w", err)
	}

	log.Printf("command %s has started as %d", scr, cmd.Pid())

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
