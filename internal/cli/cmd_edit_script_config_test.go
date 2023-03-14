package cli_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/suite"
)

type CmdEditScriptConfigTestSuite struct {
	CmdTestSuite
}

func (suite *CmdEditScriptConfigTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("edit-config", cli.NewEditScriptConfig)

	suite.EnsureScriptConfig("ns", "s", `description = "aaa"`)
}

func (suite *CmdEditScriptConfigTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "xx", "x"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	reader, err := os.Open(paths.ConfigNamespaceScriptConfig("xx", "x"))
	suite.NoError(err)

	defer reader.Close()

	_, err = config.Parse(reader)
	suite.NoError(err)

	stat, err := reader.Stat()

	suite.NoError(err)
	suite.Greater(stat.Size(), int64(0))
}

func (suite *CmdEditScriptConfigTestSuite) TestExistingFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "ns", "s"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	data, err := os.ReadFile(paths.ConfigNamespaceScriptConfig("ns", "s"))

	suite.NoError(err)
	suite.Equal(`description = "aaa"`, string(data))
}

func (suite *CmdEditScriptConfigTestSuite) TestEditorFailed() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "false", "ns", "s"})

	suite.ErrorContains(err, "command exited with 1")
	suite.NotEmpty(ctx.StdoutLines())
	suite.NotEmpty(ctx.StderrLines())

	data, err := os.ReadFile(paths.ConfigNamespaceScriptConfig("ns", "s"))

	suite.NoError(err)
	suite.Equal(`description = "aaa"`, string(data))
}

func TestCmdEditScriptConfig(t *testing.T) {
	if _, err := exec.LookPath("true"); err != nil {
		t.Skipf("cannot find out true in PATH: %v", err)
	}

	if _, err := exec.LookPath("false"); err != nil {
		t.Skipf("cannot find out false in PATH: %v", err)
	}

	suite.Run(t, &CmdEditScriptConfigTestSuite{})
}
