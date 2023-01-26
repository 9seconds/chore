package cli_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/stretchr/testify/suite"
)

type CmdEditScriptTestSuite struct {
	CmdTestSuite
}

func (suite *CmdEditScriptTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("edit-script", cli.NewEditScript)

	suite.EnsureScript("ns", "s", "echo 1")
}

func (suite *CmdEditScriptTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "xx", "x"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	stat, err := os.Stat(suite.ConfigScriptPath("xx", "x"))

	suite.NoError(err)
	suite.Greater(stat.Size(), int64(0))
}

func (suite *CmdEditScriptTestSuite) TestExistingFile() {
	before, err := os.ReadFile(suite.ConfigScriptPath("ns", "s"))
	suite.NoError(err)

	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "ns", "s"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	after, err := os.ReadFile(suite.ConfigScriptPath("ns", "s"))
	suite.NoError(err)

	suite.NoError(err)
	suite.Equal(before, after)
}

func (suite *CmdEditScriptTestSuite) TestEditorFailed() {
	before, err := os.ReadFile(suite.ConfigScriptPath("ns", "s"))
	suite.NoError(err)

	ctx, err := suite.ExecuteCommand([]string{"-e", "false", "ns", "s"})

	suite.ErrorContains(err, "command exited with 1")
	suite.NotEmpty(ctx.StdoutLines())
	suite.NotEmpty(ctx.StderrLines())

	after, err := os.ReadFile(suite.ConfigScriptPath("ns", "s"))
	suite.NoError(err)

	suite.NoError(err)
	suite.Equal(before, after)
}

func TestCmdEditScript(t *testing.T) {
	if _, err := exec.LookPath("true"); err != nil {
		t.Skipf("cannot find out true in PATH: %v", err)
	}

	if _, err := exec.LookPath("false"); err != nil {
		t.Skipf("cannot find out false in PATH: %v", err)
	}

	suite.Run(t, &CmdEditScriptTestSuite{})
}
