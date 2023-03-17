package validators

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

var (
	ErrNamespaceIsNotDirectory = errors.New("namespace is not a directory")
	ErrNamespaceInvalid        = errors.New("namespace is invalid")
	ErrScriptInvalid           = errors.New("script is invalid")

	asciiNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func ArgumentOptional(index int, callback cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) <= index {
			return nil
		}

		return callback(cmd, args)
	}
}

func ASCIIName(index int, err error) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if !asciiNameRegexp.MatchString(args[index]) {
			return err
		}

		return nil
	}
}

func Namespace(index int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		namespace, exists := script.ExtractRealNamespace(args[index])
		if !exists {
			return ErrNamespaceInvalid
		}

		stat, err := os.Stat(paths.ConfigNamespace(namespace))
		if err != nil {
			return fmt.Errorf("invalid namespace: %w", err)
		}

		if !stat.IsDir() {
			return ErrNamespaceIsNotDirectory
		}

		return nil
	}
}

func Script(nsIndex, scrIndex int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		namespace, exists := script.ExtractRealNamespace(args[nsIndex])
		if !exists {
			return ErrNamespaceInvalid
		}

		scr := &script.Script{
			Namespace:  namespace,
			Executable: args[scrIndex],
		}

		if err := script.ValidateScript(scr.Path()); err != nil {
			return ErrScriptInvalid
		}

		return nil
	}
}
