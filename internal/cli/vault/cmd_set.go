package vault

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/cli/base"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var ErrStdinIsNotTerminal = errors.New("stdin is not connected to a valid terminal")

func NewSet() *cobra.Command {
	return &cobra.Command{
		Use:                   "set namespace key [value]",
		Aliases:               []string{"s"},
		Short:                 "Set a vault secret",
		ValidArgsFunction:     completeSecretKey,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.RangeArgs(2, 3), //nolint: gomnd
			validators.Namespace(0),
		),
		Run: base.Main(main(func(cmd *cobra.Command, vlt vault.Vault, args []string) (bool, error) {
			var (
				value string
				err   error
			)

			if len(args) == 1 {
				value, err = mainSetReadFromTerminal(cmd)
			} else {
				value = args[1]
			}

			if err != nil {
				return true, fmt.Errorf("cannot set value: %w", err)
			}

			vlt.Set(args[0], value)

			return true, nil
		})),
	}
}

func mainSetReadFromTerminal(cmd *cobra.Command) (string, error) {
	descr := int(os.Stdin.Fd())

	if !term.IsTerminal(descr) {
		return "", ErrStdinIsNotTerminal
	}

	for {
		cmd.Print("Enter value: ")

		value, err := term.ReadPassword(descr)

		cmd.Println()

		if err != nil {
			return "", fmt.Errorf("cannot read value: %w", err)
		}

		cmd.Print("Repeat value: ")

		repeat, err := term.ReadPassword(descr)

		cmd.Println()

		if err != nil {
			return "", fmt.Errorf("cannot read value: %w", err)
		}

		if bytes.Equal(value, repeat) {
			return string(value), nil
		}

		cmd.Println("Value mismatch. Please try again.")
	}
}
