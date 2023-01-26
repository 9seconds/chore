package cli_test

import (
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CmdTestSuite struct {
	suite.Suite
	testlib.CobraTestSuite

	subcommand string
}

func (suite *CmdTestSuite) Setup(subcommand string, makeCommand func() *cobra.Command) {
	suite.CobraTestSuite.Setup(suite.T(), func() *cobra.Command {
		root := cli.NewRoot("version")

		root.AddCommand(makeCommand())

		return root
	})

	suite.subcommand = subcommand
}

func (suite *CmdTestSuite) ExecuteCommand(args []string) (*testlib.CobraCommandContext, error) {
	args = append([]string{suite.subcommand}, args...)

	return suite.CobraTestSuite.ExecuteCommand(args)
}
