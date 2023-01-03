package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdRun struct {
	Timeout cli.Timeout `short:"t" help:"Limit execution time."`

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
