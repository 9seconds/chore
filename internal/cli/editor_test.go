package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/stretchr/testify/suite"
)

type EditorTestSuite struct {
	suite.Suite

	param cli.Editor
}

func (suite *EditorTestSuite) SetupTest() {
	suite.param = cli.Editor("")
}

func (suite *EditorTestSuite) SetEnv(key, value string) {
	suite.T().Setenv(key, value)
}

func (suite *EditorTestSuite) TestParseUnknown() {
	suite.ErrorContains(
		suite.param.UnmarshalText([]byte("-1-1-1-1-1")),
		"cannot detect given editor")
}

func (suite *EditorTestSuite) TestParseOk() {
	suite.NoError(suite.param.UnmarshalText([]byte("go")))

	executable, err := suite.param.Value()
	suite.NoError(err)
	suite.Contains(executable, "go")
}

func (suite *EditorTestSuite) TestVisualPreferences() {
	suite.SetEnv("VISUAL", "code")
	suite.SetEnv("EDITOR", "vim")
	suite.SetEnv("PATH", "")

	executable, err := suite.param.Value()
	suite.NoError(err)
	suite.Contains(executable, "code")
}

func (suite *EditorTestSuite) TestEditorPreferences() {
	suite.SetEnv("VISUAL", "")
	suite.SetEnv("EDITOR", "vim")
	suite.SetEnv("PATH", "")

	executable, err := suite.param.Value()
	suite.NoError(err)
	suite.Contains(executable, "vim")
}

func (suite *EditorTestSuite) TestCannotFind() {
	suite.SetEnv("VISUAL", "")
	suite.SetEnv("EDITOR", "")
	suite.SetEnv("PATH", "")

	_, err := suite.param.Value()
	suite.ErrorContains(err, "cannot find out")
}

func TestEditor(t *testing.T) {
	suite.Run(t, &EditorTestSuite{})
}
