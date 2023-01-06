package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

var (
	ErrNamespaceIsNotDirectory = errors.New("namespace is not a directory")
	ErrNamespaceInvalid        = errors.New("namespace is invalid")
	ErrScriptInvalid           = errors.New("script is invalid")

	ValidNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

func argumentOptional(index int, callback cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) <= index {
			return nil
		}

		return callback(cmd, args)
	}
}

func validScriptName(index int, err error) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if !ValidNameRegexp.MatchString(args[index]) {
			return err
		}

		return nil
	}
}

func validNamespace(index int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		stat, err := os.Stat(env.PathConfigNamespace(args[index]))
		if err != nil {
			return fmt.Errorf("invalid namespace: %w", err)
		}

		if !stat.IsDir() {
			return ErrNamespaceIsNotDirectory
		}

		return nil
	}
}

func validScript(nsIndex, scrIndex int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		scr := &script.Script{
			Namespace:  args[nsIndex],
			Executable: args[scrIndex],
		}

		return script.ValidateScript(scr.Path())
	}
}
