package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdGCTestSuite struct {
	CmdTestSuite
}

func (suite *CmdGCTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("gc", cli.NewGC)

	suite.EnsureDir(suite.ConfigNamespacePath("xx"))
	suite.EnsureFile(suite.ConfigScriptPath("xy", "y"), "11", 0o600)
	suite.EnsureFile(suite.ConfigScriptPath("xx", "z"), "11", 0o600)
}

func (suite *CmdGCTestSuite) TestDryRun() {
	for _, testValue := range []string{"-n", "--dry-run"} {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			ctx, err := suite.ExecuteCommand([]string{testValue})

			assert.NoError(t, err)
			assert.DirExists(t, suite.ConfigNamespacePath("xx"))
			assert.FileExists(t, suite.ConfigScriptPath("xy", "y"))
			assert.FileExists(t, suite.ConfigScriptPath("xx", "z"))
			assert.Empty(t, ctx.StderrLines())

			lines := ctx.StdoutLines()

			assert.NotEmpty(t, lines)
			assert.Len(t, lines, 2)
			assert.Contains(t, lines, suite.ConfigNamespacePath("xx"))
			assert.Contains(t, lines, suite.ConfigNamespacePath("xy"))
		})
	}
}

func (suite *CmdGCTestSuite) TestRun() {
	ctx, err := suite.ExecuteCommand([]string{})

	suite.NoError(err)
	suite.NoDirExists(suite.ConfigNamespacePath("xx"))
	suite.NoDirExists(suite.ConfigNamespacePath("xy"))
	suite.Empty(ctx.StderrLines())
	suite.Empty(ctx.StdoutLines())
}

func TestCmdGC(t *testing.T) {
	suite.Run(t, &CmdGCTestSuite{})
}
