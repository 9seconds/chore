package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/suite"
)

type CmdRemoveTestSuite struct {
	CmdTestSuite
}

func (suite *CmdRemoveTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("rm", cli.NewRemove)

	suite.EnsureFile(paths.ConfigNamespaceScript("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.ConfigNamespaceScriptConfig("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.CacheNamespaceScript("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.DataNamespaceScript("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.DataNamespaceScript("xx", "y"), "", 0o666)
	suite.EnsureFile(paths.StateNamespaceScript("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.StateNamespaceScript("xx", "x"), "", 0o666)
	suite.EnsureFile(paths.StateNamespaceScript("xx", "z"), "", 0o666)
}

func (suite *CmdRemoveTestSuite) TestDryRun() {
	ctx, err := suite.ExecuteCommand("-n", "xx", "x", "y", "z")

	suite.NoError(err)
	suite.Empty(ctx.StderrLines())
	suite.Contains(ctx.StdoutLines(), paths.ConfigNamespaceScript("xx", "x"))
	suite.NotContains(ctx.StdoutLines(), paths.ConfigNamespaceScript("xx", "y"))
}

func (suite *CmdRemoveTestSuite) TestRun() {
	ctx, err := suite.ExecuteCommand("xx", "x", "y", "z")

	suite.NoError(err)
	suite.Empty(ctx.StderrLines())
	suite.Empty(ctx.StdoutLines())
	suite.NoFileExists(paths.ConfigNamespaceScript("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.ConfigNamespaceScriptConfig("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.CacheNamespaceScript("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.DataNamespaceScript("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.DataNamespaceScript("xx", "y"), "", 0o666)
	suite.NoFileExists(paths.StateNamespaceScript("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.StateNamespaceScript("xx", "x"), "", 0o666)
	suite.NoFileExists(paths.StateNamespaceScript("xx", "z"), "", 0o666)
}

func TestCmdRemove(t *testing.T) {
	suite.Run(t, &CmdRemoveTestSuite{})
}
