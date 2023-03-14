package cli_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/suite"
)

type CmdEditConfigTestSuite struct {
	CmdTestSuite
}

func (suite *CmdEditConfigTestSuite) SetupTest() {
	suite.CmdTestSuite.Setup("edit-config", cli.NewEditConfig)
}

func (suite *CmdEditConfigTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand([]string{"-e", "true"})
	suite.NoError(err)
	suite.Empty(ctx.StdoutLines())
	suite.Empty(ctx.StderrLines())

	reader, err := os.Open(paths.AppConfigPath())
	suite.NoError(err)

	defer reader.Close()

	conf, err := config.ReadConfig(reader)
	suite.NoError(err)

	suite.Empty(conf.Vault)
}

func TestCmdEditConfig(t *testing.T) {
	suite.Run(t, &CmdEditConfigTestSuite{})
}
