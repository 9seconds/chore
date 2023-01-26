package cli_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/config"
	"github.com/stretchr/testify/suite"
)

type CmdEditConfigTestSuite struct {
	CmdTestSuite
}

func (suite *CmdEditConfigTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("edit-config", cli.NewEditConfig)

	suite.EnsureScriptConfig("ns", "s", `description = "aaa"`)
}

func (suite *CmdEditConfigTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "xx", "x"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	reader, err := os.Open(suite.ConfigScriptConfigPath("xx", "x"))
	suite.NoError(err)

	defer reader.Close()

	_, err = config.Parse(reader)
	suite.NoError(err)

	stat, err := reader.Stat()

	suite.NoError(err)
	suite.Greater(stat.Size(), int64(0))
}

func (suite *CmdEditConfigTestSuite) TestExistingFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true", "ns", "s"})

	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	data, err := os.ReadFile(suite.ConfigScriptConfigPath("ns", "s"))

	suite.NoError(err)
	suite.Equal(`description = "aaa"`, string(data))
}

func (suite *CmdEditConfigTestSuite) TestEditorFailed() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "false", "ns", "s"})

	suite.ErrorContains(err, "command exited with 1")
	suite.NotEmpty(ctx.StdoutLines())
	suite.NotEmpty(ctx.StderrLines())

	data, err := os.ReadFile(suite.ConfigScriptConfigPath("ns", "s"))

	suite.NoError(err)
	suite.Equal(`description = "aaa"`, string(data))
}

func TestCmdEditConfig(t *testing.T) {
	if _, err := exec.LookPath("true"); err != nil {
		t.Skipf("cannot find out true in PATH: %v", err)
	}

	if _, err := exec.LookPath("false"); err != nil {
		t.Skipf("cannot find out false in PATH: %v", err)
	}

	suite.Run(t, &CmdEditConfigTestSuite{})
}
