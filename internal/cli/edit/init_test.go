package edit_test

import (
	"os/exec"

	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type EditTestSuite struct {
	suite.Suite

	testlib.CobraTestSuite
}

func (suite *EditTestSuite) SetupSuite() {
	if _, err := exec.LookPath("true"); err != nil {
		suite.T().Skipf("cannot find true in a PATH: %v", err)
	}
}

func (suite *EditTestSuite) Setup(cliCommand string, makeCommand func() *cobra.Command) {
	suite.CobraTestSuite.Setup(suite.T(), cliCommand, func() *cobra.Command {
		cmd := &cobra.Command{}

		suite.T().Setenv("VISUAL", "true")

		var editorFlag edit.FlagEditor

		cmd.PersistentFlags().VarP(&editorFlag, "editor", "e", "editor to use")
		cmd.AddCommand(makeCommand())

		return cmd
	})
}
