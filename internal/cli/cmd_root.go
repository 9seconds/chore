package cli

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/9seconds/chore/internal/env"
	"github.com/spf13/cobra"
)

func NewRoot(version string) *cobra.Command {
	var isDebug bool

	root := &cobra.Command{
		Use:     "chore",
		Short:   "A sometimes better management for your homebrew scripts.",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// env.Debug is a global flag propagated from a callee.
			// it forcefully enabled debugging for a script.
			isDebug = isDebug || os.Getenv(env.Debug) == env.DebugEnabled

			if isDebug {
				if err := os.Setenv(env.Debug, env.DebugEnabled); err != nil {
					return fmt.Errorf("cannot set debug environment variable: %w", err)
				}
			} else {
				log.SetOutput(io.Discard)
			}

			return nil
		},
		TraverseChildren: true,
	}

	flags := root.Flags()

	flags.BoolVarP(&isDebug, "debug", "d", false, "run in debug mode")

	return root
}
