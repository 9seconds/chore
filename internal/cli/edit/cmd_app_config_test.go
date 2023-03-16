package edit_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/suite"
)

type AppConfigTestSuite struct {
	EditTestSuite
}

func (suite *AppConfigTestSuite) SetupTest() {
	suite.EditTestSuite.Setup("app-config", edit.NewAppConfig)
}

func (suite *AppConfigTestSuite) TestNewFile() {
	ctx, err := suite.ExecuteCommand()
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

func TestAppConfig(t *testing.T) {
	suite.Run(t, &AppConfigTestSuite{})
}
