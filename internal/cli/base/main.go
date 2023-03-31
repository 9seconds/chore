package base

import (
	"errors"
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/paths"
	"github.com/spf13/cobra"
)

var ExitFunc = os.Exit

type ErrExit struct {
	Code int
}

func (e ErrExit) Error() string {
	return fmt.Sprintf("exited with %d code", e.Code)
}

func Main(callback func(cmd *cobra.Command, args []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		err := callback(cmd, args)
		exitErr := ErrExit{}

		paths.TempDirCleanup()

		switch {
		case errors.As(err, &exitErr):
			ExitFunc(exitErr.Code)
		case err != nil:
			cmd.PrintErrln(err)
			ExitFunc(1)
		}
	}
}
