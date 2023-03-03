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
		Use:     "run namespace script [options] [--] [args]",
		Aliases: []string{"r"},
		Short:   "Run chore script",
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(2), //nolint: gomnd
			validASCIIName(0, ErrNamespaceInvalid),
			validASCIIName(1, ErrScriptInvalid),
			validNamespace(0),
			validScript(0, 1),
		),
		Run:                   mainRun,
		ValidArgsFunction:     completeRun,
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}

	// workaround to get rid of -h/--help flag
	cmd.InitDefaultHelpFlag()

	if err := cmd.Flags().MarkHidden("help"); err != nil {
		panic(err)
	}

	return cmd
}

func mainRun(cmd *cobra.Command, args []string) {
	exitCode, err := mainRunWrapper(cmd, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command has failed: %v\n", err)
		commands.Exit(1)
	}

	if exitCode != 0 {
		commands.Exit(exitCode)
	}
}

func mainRunWrapper(cmd *cobra.Command, args []string) (int, error) {
	listDelimiter, err := cmd.Root().Flags().GetString("list-delimiter")
	if err != nil {
		return 0, fmt.Errorf("cannot get a value of 'list-delimiter' flag: %w", err)
	}

	ctx := cmd.Context()

	scr, err := script.New(args[0], args[1])
	if err != nil {
		return 0, fmt.Errorf("cannot initialize script: %w", err)
	}

	if err := scr.EnsureDirs(); err != nil {
		return 0, fmt.Errorf("cannot initialize script directories: %w", err)
	}

	parsedArgs, err := argparse.Parse(args[2:], listDelimiter)
	if err != nil {
		return 0, fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := parsedArgs.Validate(ctx, scr.Config.Flags, scr.Config.Parameters); err != nil {
		return 0, fmt.Errorf("cannot validate arguments: %w", err)
	}

	environ := scr.Environ(ctx, parsedArgs)

	for _, v := range environ {
		log.Printf("env: %s", v)
	}

	runCmd := commands.New(
		scr.Path(),
		parsedArgs.Positional,
		environ,
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
