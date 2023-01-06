package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run [flags] namespace script [options] [--] [args]",
		Aliases: []string{"r"},
		Short:   "Run chore script",
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(2), //nolint: gomnd
			validScriptName(0, ErrNamespaceInvalid),
			validScriptName(1, ErrScriptInvalid),
			validNamespace(0),
			validScript(0, 1),
		),
		Run:               mainRun,
		ValidArgsFunction: completeRun,
	}

	return cmd
}

func mainRun(cmd *cobra.Command, args []string) {
	exitCode, err := mainRunWrapper(cmd, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command has failed: %v\n", err)

		exitCode = 1
	}

	os.Exit(exitCode)
}

func mainRunWrapper(cmd *cobra.Command, args []string) (int, error) {
	ctx := cmd.Context()
	scr := &script.Script{
		Namespace:  args[0],
		Executable: args[1],
	}

	if err := scr.Init(); err != nil {
		return 0, fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	conf := scr.Config()

	parsedArgs, err := argparse.Parse(args[2:])
	if err != nil {
		return 0, fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := parsedArgs.Validate(ctx, conf.Flags, conf.Parameters); err != nil {
		return 0, fmt.Errorf("cannot validate arguments: %w", err)
	}

	environ := scr.Environ(ctx, parsedArgs)

	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	runCmd := commands.NewOS(
		scr,
		environ,
		parsedArgs.Positional,
		os.Stdin,
		os.Stdout,
		os.Stderr)

	if err := runCmd.Start(ctx); err != nil {
		return 0, fmt.Errorf("cannot start command: %w", err)
	}

	log.Printf("command %s has started as %d", scr, runCmd.Pid())

	result := runCmd.Wait()

	log.Printf("command %d exit with exit code %d", runCmd.Pid(), result.ExitCode)
	log.Printf(
		"command %d times: user=%v, sys=%v, real=%v",
		runCmd.Pid(),
		result.UserTime,
		result.SystemTime,
		result.ElapsedTime)

	return result.ExitCode, nil
}
