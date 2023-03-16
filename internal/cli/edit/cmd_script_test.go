package edit_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/suite"
)

type ScriptTestSuite struct {
	EditTestSuite
}

func (suite *ScriptTestSuite) SetupTest() {
	suite.EditTestSuite.Setup("script", edit.NewScript)
	suite.EnsureScript("ns", "s", "echo 1")
}

func (suite *ScriptTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand("xx", "x")

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	stat, err := os.Stat(paths.ConfigNamespaceScript("xx", "x"))

	suite.NoError(err)
	suite.Greater(stat.Size(), int64(0))
}

func (suite *ScriptTestSuite) TestExistingFile() {
	before, err := os.ReadFile(paths.ConfigNamespaceScript("ns", "s"))
	suite.NoError(err)

	ctx, err := suite.ExecuteCommand("ns", "s")

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	after, err := os.ReadFile(paths.ConfigNamespaceScript("ns", "s"))
	suite.NoError(err)

	suite.NoError(err)
	suite.Equal(before, after)
}

func TestScript(t *testing.T) {
	suite.Run(t, &ScriptTestSuite{})
}
