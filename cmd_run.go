package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
	"github.com/gofrs/flock"
)

const fileLockPollPeriod = 100 * time.Millisecond

type CliCmdRun struct {
	Timeout cli.Timeout `short:"t" help:"Limit execution time."`
	Lock    bool        `short:"l" help:"Acquire exclusive file lock on a script."`

	Namespace cli.Namespace `arg:"" help:"Prefix of the script namespace."`
	Script    string        `arg:"" help:"Prefix of the script name."`
	Args      []string      `arg:"" optional:"" passthrough:"" help:"Script arguments to use."`
}

func (c *CliCmdRun) Run(appCtx cli.Context) error {
	var (
		ctx    context.Context = appCtx
		cancel context.CancelFunc
	)

	if c.Timeout.Value() != 0 {
		ctx, cancel = context.WithTimeout(ctx, c.Timeout.Value())
		defer cancel()
	}

	scr, err := script.FindScript(c.Namespace.Value(), c.Script)
	if err != nil {
		return fmt.Errorf("cannot find out script: %w", err)
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	conf := scr.Config()

	args, err := argparse.Parse(c.Args)
	if err != nil {
		return fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := args.Validate(ctx, conf.Flags, conf.Parameters); err != nil {
		return fmt.Errorf("cannot validate arguments: %w", err)
	}

	environ := scr.Environ(ctx, args)

	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	cmd := commands.NewOS(
		scr,
		environ,
		args.Positional,
		os.Stdin,
		os.Stdout,
		os.Stderr)

	lock, err := c.lockScript(ctx, scr.Path())
	if err != nil {
		return err
	}

	defer lock.Unlock() //nolint: errcheck

	if err := cmd.Start(ctx); err != nil {
		return fmt.Errorf("cannot start command: %w", err)
	}

	log.Printf("command %s has started as %d", scr, cmd.Pid())

	result := cmd.Wait()
	if !result.Ok() {
		return fmt.Errorf("cannot correctly finish command: %w", result)
	}

	log.Printf("command %d exit with exit code %d", cmd.Pid(), result.ExitCode)
	log.Printf(
		"command %d times: user=%v, sys=%v, real=%v",
		cmd.Pid(),
		result.UserTime,
		result.SystemTime,
		result.ElapsedTime)

	return result
}

func (c *CliCmdRun) lockScript(ctx context.Context, path string) (*flock.Flock, error) {
	var (
		acquired bool
		err      error
	)

	lock := flock.New(path)

	if c.Lock {
		acquired, err = lock.TryLockContext(ctx, fileLockPollPeriod)
	} else {
		acquired, err = lock.TryRLockContext(ctx, fileLockPollPeriod)
	}

	switch {
	case err != nil:
		return nil, fmt.Errorf("cannot acquire lock on %s: %w", path, err)
	case !acquired:
		return nil, fmt.Errorf("cannot acquire lock on %s", path)
	}

	return lock, nil
}
