package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli/base"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/commands"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/env"
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
			validators.Script(0, 1),
		),
		Run:                   base.Main(mainRun),
		ValidArgsFunction:     completeRun,
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
	}

	// workaround to get rid of -h/--help flag
	cmd.InitDefaultHelpFlag()

	if err := cmd.Flags().MarkHidden("help"); err != nil {
		panic(err.Error())
	}

	return cmd
}

func mainRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	namespace, _ := script.ExtractRealNamespace(args[0])

	conf, err := config.Get()
	if err != nil {
		return fmt.Errorf("cannot open application config: %w", err)
	}

	scr, err := script.New(namespace, args[1])
	if err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	if err := scr.EnsureDirs(); err != nil {
		return fmt.Errorf("cannot initialize script directories: %w", err)
	}

	parsedArgs, err := argparse.Parse(args[2:])
	if err != nil {
		return fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := parsedArgs.Validate(ctx, scr.Config.Flags, scr.Config.Parameters); err != nil {
		return fmt.Errorf("cannot validate arguments: %w", err)
	}

	confEnviron := conf.Environ(namespace)
	for _, v := range confEnviron {
		log.Printf("config env: %s", v)
	}

	scriptEnviron := scr.Environ(ctx, parsedArgs)
	for _, v := range scriptEnviron {
		log.Printf("script env: %s", v)
	}

	environ := env.Environ()
	environ = append(environ, confEnviron...)
	environ = append(environ, scriptEnviron...)

	runCmd := commands.New(
		scr.Path(),
		parsedArgs.Positional,
		environ,
		os.Stdin,
		os.Stdout,
		os.Stderr)

	if err := runCmd.Start(ctx); err != nil {
		return fmt.Errorf("cannot start command: %w", err)
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

	return base.ErrExit{
		Code: result.ExitCode,
	}
}
