package cli_test

import (
	"path/filepath"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CmdRenameTestSuite struct {
	CmdTestSuite
}

func (suite *CmdRenameTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("rename", cli.NewRename)

	suite.EnsureScript("ns", "s", "echo 1")

	scr, err := script.New("ns", "s")
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), scr.EnsureDirs())

	suite.EnsureFile(filepath.Join(scr.CachePath(), "a"), "xx", 0o666)
}

func (suite *CmdRenameTestSuite) TestCannotMoveUnknownScript() {
	_, err := suite.ExecuteCommand("ns", "v", "x")
	suite.ErrorIs(err, validators.ErrScriptInvalid)
}

func (suite *CmdRenameTestSuite) TestHasOffendingDirectory() {
	suite.ExitMock(1).Once()

	suite.EnsureDir(paths.CacheNamespaceScript("ns", "x"))

	ctx, err := suite.ExecuteCommand("ns", "s", "x")
	suite.NoError(err)
	suite.Contains(ctx.Stderr.String(), "it prevents renaming")
}

func (suite *CmdRenameTestSuite) TestHasOffendingDirectoryForce() {
	suite.EnsureDir(paths.CacheNamespaceScript("ns", "x"))

	_, err := suite.ExecuteCommand("-f", "ns", "s", "x")
	suite.NoError(err)

	suite.NoFileExists(filepath.Join(paths.CacheNamespaceScript("ns", "s"), "a"))
	suite.FileExists(filepath.Join(paths.CacheNamespaceScript("ns", "x"), "a"))
}

func (suite *CmdRenameTestSuite) TestOk() {
	_, err := suite.ExecuteCommand("ns", "s", "x")
	suite.NoError(err)

	suite.NoFileExists(filepath.Join(paths.CacheNamespaceScript("ns", "s"), "a"))
	suite.FileExists(filepath.Join(paths.CacheNamespaceScript("ns", "x"), "a"))
}

func TestCmdRename(t *testing.T) {
	suite.Run(t, &CmdRenameTestSuite{})
}
