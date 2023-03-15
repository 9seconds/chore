package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdGCTestSuite struct {
	CmdTestSuite
}

func (suite *CmdGCTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("gc", cli.NewGC)

	suite.EnsureDir(paths.ConfigNamespace("xx"))
	suite.EnsureFile(paths.ConfigNamespaceScript("xy", "y"), "11", 0o600)
	suite.EnsureFile(paths.ConfigNamespaceScript("xx", "z"), "11", 0o600)
}

func (suite *CmdGCTestSuite) TestDryRun() {
	for _, testValue := range []string{"-n", "--dry-run"} {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			ctx, err := suite.ExecuteCommand(testValue)

			assert.NoError(t, err)
			assert.DirExists(t, paths.ConfigNamespace("xx"))
			assert.FileExists(t, paths.ConfigNamespaceScript("xy", "y"))
			assert.FileExists(t, paths.ConfigNamespaceScript("xx", "z"))
			assert.Empty(t, ctx.StderrLines())

			lines := ctx.StdoutLines()

			assert.NotEmpty(t, lines)
			assert.Len(t, lines, 2)
			assert.Contains(t, lines, paths.ConfigNamespace("xx"))
			assert.Contains(t, lines, paths.ConfigNamespace("xy"))
		})
	}
}

func (suite *CmdGCTestSuite) TestRun() {
	ctx, err := suite.ExecuteCommand()

	suite.NoError(err)
	suite.NoDirExists(paths.ConfigNamespace("xx"))
	suite.NoDirExists(paths.ConfigNamespace("xy"))
	suite.Empty(ctx.StderrLines())
	suite.Empty(ctx.StdoutLines())
}

func TestCmdGC(t *testing.T) {
	suite.Run(t, &CmdGCTestSuite{})
}
