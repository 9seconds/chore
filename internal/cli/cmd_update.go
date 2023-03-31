package cli

import (
	"fmt"

	"github.com/9seconds/chore/internal/cli/base"
	"github.com/9seconds/chore/internal/update"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

func NewUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update",
		Short:             "Check, verify and update binary in place.",
		Args:              cobra.NoArgs,
		Run:               base.Main(mainSelfUpdate),
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	flags := cmd.Flags()
	flags.BoolP("check", "c", false, "only check if update is required")
	flags.BoolP("unstable", "u", false, "check also unstable versions")

	return cmd
}

func mainSelfUpdate(cmd *cobra.Command, _ []string) error {
	checkFlag, _ := cmd.Flags().GetBool("check")
	unstableFlag, _ := cmd.Flags().GetBool("unstable")

	release, err := update.GetLatestRelease(cmd.Context(), unstableFlag)
	if err != nil {
		return err
	}

	version := cmd.Root().Version

	if version == release.Version {
		cmd.Println("Nothing to update")

		return nil
	}

	cmd.Printf("Update %s -> %s\n", version, release.Version)

	if checkFlag {
		return nil
	}

	binary, err := update.Extract(
		cmd.Context(),
		release.ArchiveURL,
		release.SignatureURL)
	if err != nil {
		return fmt.Errorf("cannot get binary: %w", err)
	}

	if err := selfupdate.RollbackError(selfupdate.Apply(binary, selfupdate.Options{})); err != nil {
		return fmt.Errorf("cannot perform update: %w", err)
	}

	cmd.Println("Updated!")

	return nil
}
