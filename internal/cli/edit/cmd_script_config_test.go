package edit_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/suite"
)

type ScriptConfigTest struct {
	EditTestSuite
}

func (suite *ScriptConfigTest) SetupTest() {
	suite.EditTestSuite.Setup("script-config", edit.NewScriptConfig)
	suite.EnsureScript("ns", "z", "echo 1")
	suite.EnsureScript("ns", "s", "echo 2")
	suite.EnsureScriptConfig("ns", "s", `description = "aaa"`)
}

func (suite *ScriptConfigTest) TestUnknownNamespace() {
	_, err := suite.ExecuteCommand("xx", "x")
	suite.ErrorIs(err, validators.ErrScriptInvalid)
}

func (suite *ScriptConfigTest) TestUnknownScript() {
	_, err := suite.ExecuteCommand("ns", "x")
	suite.ErrorIs(err, validators.ErrScriptInvalid)
}

func (suite *ScriptConfigTest) TestNewFile() {
	ctx, err := suite.ExecuteCommand("ns", "z")

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	reader, err := os.Open(paths.ConfigNamespaceScriptConfig("ns", "z"))
	suite.NoError(err)

	defer reader.Close()

	_, err = config.Parse(reader)
	suite.NoError(err)

	stat, err := reader.Stat()

	suite.NoError(err)
	suite.Greater(stat.Size(), int64(0))
}

func (suite *ScriptConfigTest) TestExistingFile() {
	ctx, err := suite.ExecuteCommand("ns", "s")

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	data, err := os.ReadFile(paths.ConfigNamespaceScriptConfig("ns", "s"))
	suite.NoError(err)
	suite.Equal(`description = "aaa"`, string(data))
}

func TestScriptConfig(t *testing.T) {
	suite.Run(t, &ScriptConfigTest{})
}
